package users

import (
	"github.com/StefanPahlplatz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/date_utils"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/errors"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, created_on) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, created_on FROM users WHERE id=?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer sqlStatement.Close()

	// Get the row.
	result := sqlStatement.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedOn); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

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

	// Insert the user.
	insertResult, saveErr := sqlStatement.Exec(user.FirstName, user.LastName, user.Email, user.CreatedOn)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	// Set the ID.
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(saveErr)
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer sqlStatement.Close()

	// Update the user
	_, updateErr := sqlStatement.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer sqlStatement.Close()

	// Delete the user.
	if _, err := sqlStatement.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}
