package models

/*
  @Author : Mustang Kong
*/

import (
	"golang-common-base/app/models/auth"
	"golang-common-base/app/models/email"
	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
)

func AutoMigrateTable() {
	connection.DB.Self.AutoMigrate(
		// auth
		&auth.User{},

		// sys
		&sys.User{},
		&sys.Role{},
		&sys.Menu{},
		&sys.Dept{},
		&sys.UserRole{},
		&sys.RoleMenu{},

		&email.EmailTextContent{},
	)
}
