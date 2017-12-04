package jwt

import (
	"io/ioutil"
	"testing"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/zero-os/0-stor/client/itsyouonline"
	"github.com/zero-os/0-stor/stubs"
)

func BenchmarkJWTCache(b *testing.B) {
	require := require.New(b)

	token := getTokenBench(require)
	v, err := getTestVerifier(true)
	require.NoError(err, "failed to create jwt verifier")
	verifier := v.(*Verifier)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := verifier.getScopes(token)
		require.NoError(err, "getScopes failed")
	}
}

func BenchmarkJWTWithoutCache(b *testing.B) {
	require := require.New(b)

	token := getTokenBench(require)
	v, err := getTestVerifier(true)
	require.NoError(err, "failed to create jwt verifier")
	verifier := v.(*Verifier)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _, err := verifier.checkJWTGetScopes(token)
		if err != nil {
			require.NoError(err, "getScopes failed")
		}
	}
}

func getTokenBench(require *require.Assertions) string {
	b, err := ioutil.ReadFile("../../devcert/jwt_key.pem")
	require.NoError(err)

	key, err := jwtgo.ParseECPrivateKeyFromPEM(b)
	require.NoError(err)

	iyoCl, err := stubs.NewStubIYOClient("testorg", key)

	token, err := iyoCl.CreateJWT("mynamespace", itsyouonline.Permission{
		Read:   true,
		Write:  true,
		Delete: true,
	})
	require.NoError(err)

	return token
}
