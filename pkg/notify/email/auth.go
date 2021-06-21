package email

import (
	"errors"
	"net/smtp"
)

type customAuth struct {
	username, password string
}

func CustomAuth(username, password string) smtp.Auth {
	return &customAuth{username: username, password: password}
}

func (ca *customAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	proto = "LOGIN"
	return
}

func (ca *customAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(ca.username), err
		case "Password:":
			return []byte(ca.password), err
		default:
			return nil, errors.New("unkown fromserver")
		}
	}
	return
}
