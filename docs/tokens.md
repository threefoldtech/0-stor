# 0-stor tokens

## Different tokens:
The 0-stor is used 3 type of tokens.
- ItYou.online JWT token
- Reservation token
- Data Access token

### ItYou.online JWT token
This token is used to authenticate the user that creates sends the requests.
It must contains the scope `user:name`.

### Reservation token
This token is generated by the 0-stor when a new reservation is created.  
This token is required when you want to renew a reservation or when you want to inspect the stats of a reservation.

This token contains:
- The id of the reservation
- The username of the user that created the reservation
- The size reserved
- The namespace label


### Data Access token
This token is created when a new reservation is created or when we want to authorize a user to access a namespace.  
A user need to pass this token in the requests when he wants to access the data in a namespace.

This token contains:
- The namespace label
- The username of the user we want to authorize
- The permission given to the user (read, write, delete, admin)