package users

import (
	"fmt"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/errors"
)

var (
	userDb = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
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
	current := userDb[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}
	userDb[user.Id] = user
	return nil
}
