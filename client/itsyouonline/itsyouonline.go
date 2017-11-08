package itsyouonline

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/itsyouonline/identityserver/clients/go/itsyouonline"
)

const (
	accessTokenURI = "https://itsyou.online/v1/oauth/access_token?response_type=id_token"
)

var (
	errNoPermission = errors.New("no permission")

	// ErrForbidden represents a forbidden action error
	ErrForbidden = errors.New("forbidden action")
)

// IYOClient defines the interface to manage namespaces and permissions on ItsYouOnline
type IYOClient interface {
	CreateJWT(namespace string, perm Permission) (string, error)
	CreateNamespace(namespace string) error
	DeleteNamespace(namespace string) error
	GivePermission(namespace, userID string, perm Permission) error
	RemovePermission(namespace, userID string, perm Permission) error
	GetPermission(namespace, userID string) (Permission, error)
}

// Client defines itsyouonline client which is designed to help 0-stor user.
// It is not replacement for official itsyouonline client
type Client struct {
	org       string
	clientID  string
	secret    string
	iyoClient *itsyouonline.Itsyouonline
}

// NewClient creates new client
func NewClient(org, clientID, secret string) *Client {
	return &Client{
		org:       org,
		clientID:  clientID,
		secret:    secret,
		iyoClient: itsyouonline.NewItsyouonline(),
	}
}

// CreateJWT creates itsyouonline JWT token with these scopes:
// - org.namespace.read if perm.Read is true
// - org.namespace.write if perm.Write is true
// - org.namespace.delete if perm.Delete is true
func (c *Client) CreateJWT(namespace string, perm Permission) (string, error) {
	qp := map[string]interface{}{
		"grant_type":    "client_credentials",
		"client_id":     c.clientID,
		"client_secret": c.secret,
	}

	// build scopes query
	scopes := perm.Scopes(c.org, "0stor"+"."+namespace)
	if len(scopes) == 0 {
		return "", errNoPermission
	}
	qp["scope"] = strings.Join(scopes, ",")

	// create the request
	req, err := http.NewRequest("POST", accessTokenURI, nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = buildQueryString(req, qp)

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read response
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get access token, response code = %v", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err

}

// CreateNamespace creates namespace as itsyouonline organization
// It creates these organizations:
// - org.0stor.namespace.read
// - org.0stor.namespace.write
// - org.0stor.namespace.write
func (c *Client) CreateNamespace(namespace string) error {
	_, _, _, err := c.iyoClient.LoginWithClientCredentials(c.clientID, c.secret)
	if err != nil {
		return fmt.Errorf("login failed:%v", err)
	}

	// create namespace org
	namespaceID := c.org + "." + "0stor"
	org := itsyouonline.Organization{
		Globalid: namespaceID,
	}
	_, resp, err := c.iyoClient.Organizations.CreateNewSubOrganization(c.org, org, nil, nil)
	if err != nil {
		return fmt.Errorf("code=%v, err=%v", resp.StatusCode, err)
	}

	// create 0stor suborganization

	org = itsyouonline.Organization{
		Globalid: namespaceID + "." + namespace,
	}
	_, resp, err = c.iyoClient.Organizations.CreateNewSubOrganization(namespaceID, org, nil, nil)

	if err != nil {
		return fmt.Errorf("code=%v, err=%v", resp.StatusCode, err)
	}

	namespaceID = namespaceID + "." + namespace

	// create permission org
	perm := Permission{
		Read:   true,
		Delete: true,
		Write:  true,
	}
	for _, perm := range perm.perms() {
		org := itsyouonline.Organization{
			Globalid: namespaceID + "." + perm,
		}
		_, resp, err := c.iyoClient.Organizations.CreateNewSubOrganization(namespaceID, org, nil, nil)
		if err != nil {
			return fmt.Errorf("code=%v, err=%v", resp.StatusCode, err)
		}
	}
	return nil
}

// DeleteNamespace deletes the namespace sub organization and all of it's sub organizations
func (c *Client) DeleteNamespace(namespace string) error {
	_, _, _, err := c.iyoClient.LoginWithClientCredentials(c.clientID, c.secret)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	resp, err := c.iyoClient.Organizations.DeleteOrganization(c.namespaceID(namespace), nil, nil)
	if err != nil {
		return fmt.Errorf("deleting namespace failed: IYO returned status %+v \nwith error message: %v", resp.Status, err)
	}

	if resp.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}

	return nil
}

// GivePermission give a user some permission on a namespace
func (c *Client) GivePermission(namespace, userID string, perm Permission) error {
	_, _, _, err := c.iyoClient.LoginWithClientCredentials(c.clientID, c.secret)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	var org string
	for _, perm := range perm.perms() {
		if perm == "admin" {
			org = c.namespaceID(namespace)
		} else {
			org = c.namespaceID(namespace) + "." + perm
		}
		user := itsyouonline.OrganizationsGlobalidMembersPostReqBody{Searchstring: userID}
		_, resp, err := c.iyoClient.Organizations.AddOrganizationMember(org, user, nil, nil)
		if err != nil {
			return fmt.Errorf("give member permission failed: code=%v, err=%v", resp.StatusCode, err)
		}
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("give member permission failed: code=%v", resp.StatusCode)
		}
	}

	return nil
}

// RemovePermission remove some permission from a user on a namespace
func (c *Client) RemovePermission(namespace, userID string, perm Permission) error {
	_, _, _, err := c.iyoClient.LoginWithClientCredentials(c.clientID, c.secret)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	var org string
	for _, perm := range perm.perms() {
		if perm == "admin" {
			org = c.namespaceID(namespace)
		} else {
			org = c.namespaceID(namespace) + "." + perm
		}
		resp, err := c.iyoClient.Organizations.RemoveOrganizationMember(userID, org, nil, nil)
		if err != nil {
			return fmt.Errorf("removing permission failed: IYO returned status %+v \nwith error message: %v", resp.Status, err)
		}
	}

	return nil
}

// GetPermission retrieves the permission a user has for a namespace
// returns true for a right when user is member or invited to the namespace
func (c *Client) GetPermission(namespace, userID string) (Permission, error) {
	var (
		permission = Permission{}
		org        string
	)

	_, _, _, err := c.iyoClient.LoginWithClientCredentials(c.clientID, c.secret)
	if err != nil {
		return permission, fmt.Errorf("login failed: %v", err)
	}

	for _, perm := range []string{"read", "write", "delete", "admin"} {
		if perm == "admin" {
			org = c.namespaceID(namespace)
		} else {
			org = c.namespaceID(namespace) + "." + perm
		}

		invitations, resp, err := c.iyoClient.Organizations.GetInvitations(org, nil, nil)
		if err != nil {
			return permission, fmt.Errorf("Failed to retrieve user permission : %+v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return permission, fmt.Errorf("Failed to retrieve user permission : IYO returned status %+v", resp.Status)
		}

		members, resp, err := c.iyoClient.Organizations.GetOrganizationUsers(org, nil, nil)
		if err != nil {
			return permission, fmt.Errorf("Failed to retrieve user permission: %+v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return permission, fmt.Errorf("Failed to retrieve user permission : IYO returned status %+v", resp.Status)
		}

		switch perm {
		case "read":
			if hasPermission(userID, members.Users, invitations) {
				permission.Read = true
			}
		case "write":
			if hasPermission(userID, members.Users, invitations) {
				permission.Write = true
			}
		case "delete":
			if hasPermission(userID, members.Users, invitations) {
				permission.Delete = true
			}
		case "admin":
			if hasPermission(userID, members.Users, invitations) {
				permission.Admin = true
			}
		}
	}
	return permission, nil
}

func hasPermission(target string, members []itsyouonline.OrganizationUser, invitations []itsyouonline.JoinOrganizationInvitation) bool {
	return isMember(target, members) || isInvited(target, invitations)
}

func isMember(target string, list []itsyouonline.OrganizationUser) bool {
	for _, v := range list {
		if target == v.Username {
			return true
		}
	}
	return false
}

func isInvited(target string, invitations []itsyouonline.JoinOrganizationInvitation) bool {
	for _, invite := range invitations {
		if target == invite.User || target == invite.Emailaddress {
			return true
		}
	}
	return false
}

func (c *Client) namespaceID(namespace string) string {
	return c.org + "." + "0stor" + "." + namespace
}

func buildQueryString(req *http.Request, qs map[string]interface{}) string {
	q := req.URL.Query()

	for k, v := range qs {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	return q.Encode()
}
