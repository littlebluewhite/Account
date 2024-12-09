package user_server

import "fmt"

var NoID = fmt.Errorf("cannot get tokenTime by this id")

var NoUsername = fmt.Errorf("cannot get this username")

var WrongPassword = fmt.Errorf("wrong password")
