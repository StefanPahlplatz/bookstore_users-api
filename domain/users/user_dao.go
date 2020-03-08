package users

import (
	"fmt"
	"github.com/StefanPahlplatz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/date_utils"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/errors"
	"strings"
)

const (
	indexUniqueEmail = "users_email_uindex"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, created_on) VALUES(?, ?, ?, ?);"
)

var (
	userDb = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	result := userDb[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.CreatedOn = result.CreatedOn

	return nil
}

func (user *User) Save() *errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer sqlStatement.Close()

	user.CreatedOn = date_utils.GetNowString()

	// Insert the user
	insertResult, err := sqlStatement.Exec(user.FirstName, user.LastName, user.Email, user.CreatedOn)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	// Set the ID.
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
