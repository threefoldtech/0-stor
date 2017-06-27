package main

import (
	"github.com/zero-os/0-stor/store/goraml"
	"github.com/zero-os/0-stor/store/librairies/reservation"
	"time"
	"gopkg.in/validator.v2"
	"encoding/binary"
	"encoding/base64"
	"fmt"
	"errors"
)

type Reservation struct{
	Namespace string
	reservation.Reservation
}

func NewReservation(namespace string, admin string, size float64, period int) (*Reservation, error){
	creationDate := time.Now()
	expirationDate := creationDate.AddDate(0, 0, period)

	uuid, err := GenerateUUID(64)

	if err != nil{
		return nil, err
	}

	return &Reservation{namespace,
		reservation.Reservation{
		AdminId: admin,
		SizeReserved: size,
		SizeUsed: 0,
		ExpireAt: goraml.DateTime(expirationDate),
		Created: goraml.DateTime(creationDate),
		Updated: goraml.DateTime(creationDate),
		Id: uuid,
	}}, nil
}


func (s Reservation) Validate() error {

	return validator.Validate(s)
}

func (s Reservation) SizeRemaining() float64{
	return s.SizeReserved - s.SizeUsed

}


func (s Reservation) GetKey(config *settings) string{
	return fmt.Sprintf("%s%s_%s", config.Reservations.Namespaces.Prefix, s.Namespace, s.Id)
}

func (r Reservation) Save(db *Badger, config *settings) error{
	key := r.GetKey(config)
	return db.Set(key, r.ToBytes())
}

func (r *Reservation) Get(db *Badger, config *settings) (*Reservation, error){
	key := r.GetKey(config)
	v, err := db.Get(key)

	if err != nil{
		return nil, err
	}

	if v == nil{
		return nil, nil
	}

	r.FromBytes(v)
	return r, nil
}



func (s Reservation) ToBytes() []byte {
	/*
	-----------------------------------------------------------------
	SizeReserved| TotalSizeReserved |Size of CreationDate
	 8         |   8      |  2
	-----------------------------------------------------------------

	-----------------------------------------------------------------------
	Size of UpdateDate     |Size of ExpirationDate | Size ID | Size AdminID
 	    2                  |         2             |  2       |   2
	----------------------------------------------------------------------

	------------------------------------------------------------
	 CreationDate | UpdateDate | ExpirationDate | ID | AdminId
	------------------------------------------------------------

	*/

	adminId := s.AdminId
	aSize := int16(len(adminId))

	id := s.Id
	iSize := int16(len(id))

	created := []byte(time.Time(s.Created).Format(time.RFC3339))
	updated := []byte(time.Time(s.Updated).Format(time.RFC3339))
	expiration := []byte(time.Time(s.ExpireAt).Format(time.RFC3339))

	cSize := int16(len(created))
	uSize := int16(len(updated))
	eSize := int16(len(expiration))

	result := make([]byte, 26+cSize+uSize+eSize+aSize+iSize)

	copy(result[0:8], Float64bytes(s.SizeReserved))
	copy(result[8:16], Float64bytes(s.SizeUsed))

	binary.LittleEndian.PutUint16(result[16:18], uint16(cSize))
	binary.LittleEndian.PutUint16(result[18:20], uint16(uSize))
	binary.LittleEndian.PutUint16(result[20:22], uint16(eSize))
	binary.LittleEndian.PutUint16(result[22:24], uint16(iSize))
	binary.LittleEndian.PutUint16(result[24:26], uint16(aSize))

	//Creation Date size and date
	start := 26
	end := 26 + cSize
	copy(result[start:end], created)

	//update Date
	start2 := end
	end2 := end + uSize
	copy(result[start2:end2], updated)

	//ExpirationDate
	start3 := end2
	end3 := start3 + eSize
	copy(result[start3:end3], expiration)

	// ID
	start4 := end3
	end4 := start4 + iSize
	copy(result[start4:end4], []byte(id))

	// Admin ID
	start5 := end4
	end5 := start5 + aSize
	copy(result[start5:end5], []byte(adminId))
	return result
}

func (s *Reservation) FromBytes(data []byte) error{
	s.SizeReserved = Float64frombytes(data[0:8])
	s.SizeUsed = Float64frombytes(data[8:16])

	cSize := int16(binary.LittleEndian.Uint16(data[16:18]))
	uSize := int16(binary.LittleEndian.Uint16(data[18:20]))
	eSsize := int16(binary.LittleEndian.Uint16(data[20:22]))
	iSize := int16(binary.LittleEndian.Uint16(data[22:24]))
	aSize := int16(binary.LittleEndian.Uint16(data[24:26]))

	start := 26
	end := 26 + cSize

	cTime, err := time.Parse(time.RFC3339, string(data[start:end]))

	if err != nil{
		return err
	}

	start2 := end
	end2 := end + uSize

	uTime, err := time.Parse(time.RFC3339, string(data[start2:end2]))

	if err != nil{
		return err
	}

	start3 := end2
	end3 := end2 + eSsize

	eTime, err := time.Parse(time.RFC3339, string(data[start3:end3]))

	if err != nil{
		return err
	}

	start4 := end3
	end4 := start4 + iSize

	start5 := end4
	end5 := start5 + aSize


	s.Created = goraml.DateTime(cTime)
	s.Updated = goraml.DateTime(uTime)
	s.ExpireAt = goraml.DateTime(eTime)

	s.Id = string(data[start4:end4])
	s.AdminId = string(data[start5:end5])
	return nil
}




/*
	Token format
	-----------------------------------------------------------------------------------------------------------
	Random bytes |ReservationExpirationDateEpoch| namespace ID length| reservation ID length| namespaceID|resID
	    51           8                                2 bytes        |     2 bytes
	-----------------------------------------------------------------------------------------------------------
 */

func (s Reservation) GenerateTokenForReservation(db *Badger, namespaceID string)(string, error){
	nID := []byte(namespaceID)
	rID := []byte(s.Id)

	b := make([]byte, 63 + len(nID) + len(rID))

	r, err := GenerateRandomBytes(51)

	if err != nil{
		return "", err
	}

	copy(b[0:51], r)

	epoch := time.Time(s.ExpireAt).Unix()
	binary.LittleEndian.PutUint64(b[51:59], uint64(epoch))

	nSize := len(nID)
	rSize := len(rID)

	binary.LittleEndian.PutUint16(b[59:61], uint16(nSize))
	binary.LittleEndian.PutUint16(b[61:63], uint16(rSize))

	start := 63
	end := 63 + nSize
	copy(b[start:end], nID)

	start = end
	end = start + rSize
	copy(b[start:end], rID)

	token, err := base64.URLEncoding.EncodeToString(b), err

	if err != nil{
		return "", err
	}
	return token, nil
}

/*
| Random bytes  | expirationEpoch  |Admin|Read |Write|Delete|user|
|---------------|------------------|-----|-----|-----|------|----|
| 51 bytes      | 8 bytes          |1byte|1byte|1byte|1byte|    |


 */
func (s Reservation) GenerateDataAccessTokenForUser(user string, namespaceID string, acl ACLEntry) (string, error){
	b := make([]byte, 60 + len(namespaceID) + len(user))

	r, err := GenerateRandomBytes(51)

	if err != nil{
		return "", err
	}

	copy(b[0:51], r)
	epoch := time.Time(s.ExpireAt).Unix()
	binary.LittleEndian.PutUint64(b[51:59], uint64(epoch))

	copy(b[59:63], acl.ToBytes())
	copy(b[63:], []byte(user))

	token, err := base64.URLEncoding.EncodeToString(b), err

	if err != nil{
		return "", err
	}

	return token, nil
}

func (s *Reservation) ValidateReservationToken(token, namespaceID string) (string, error){
	bytes := []byte(token)

	if len(bytes) < 63{
		return "", errors.New("Reservation token is invalid")
	}

	namespaceSize := int16(binary.LittleEndian.Uint16(bytes[59:61]))
	reseIdSize := int16(binary.LittleEndian.Uint16(bytes[61:63]))

	if len(bytes) < 63 + int(namespaceSize) + int(reseIdSize){
		return "", errors.New("Reservation token is invalid")
	}

	now := time.Now()
	expiration := time.Unix(int64(binary.LittleEndian.Uint64(bytes[51:59])), 0)

	if now.After(expiration){
		return "", errors.New("Reservation token expired")
	}

	start := 63
	end := 63 + namespaceSize
	namespace := string(bytes[start:end])

	if namespace != namespaceID{
		return "", errors.New("Reservation token is invalid")
	}

	reservation := string(bytes[end:end+reseIdSize])

	return reservation, nil
}

func (s Reservation) ValidateDataAccessToken(acl ACLEntry, token, user string) error{
	bytes := []byte(token)
	if len(bytes) <= 63{
		return errors.New("Data access token is invalid")
	}

	now := time.Now()
	expiration := time.Unix(int64(binary.LittleEndian.Uint64(bytes[51:59])), 0)

	if now.After(expiration){
		return errors.New("Data access token expired")
	}

	tokenACL := ACLEntry{}
	tokenACL.FromBytes(bytes[59:63])

	if tokenACL.Admin != acl.Admin ||
		tokenACL.Read != acl.Read ||
		tokenACL.Write != acl.Write ||
		tokenACL.Delete != acl.Delete{

			return errors.New("Permission denied")
	}

	tokenUser := string(bytes[63:])

	if user != tokenUser{
		return errors.New("Invalid token for user")
	}

	return nil

}