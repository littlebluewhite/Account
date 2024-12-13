package user_server

import "fmt"

var NoUsername = fmt.Errorf("cannot get this username")

var WrongPassword = fmt.Errorf("wrong password")

var UsernameExist = fmt.Errorf("username already exist")

var NoToken = fmt.Errorf("no token")

var NoUser = fmt.Errorf("no user")
