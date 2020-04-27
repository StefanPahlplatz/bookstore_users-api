package users

import (
	"errors"
	"fmt"
	"github.com/StefanPahlplatz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/StefanPahlplatz/bookstore_users-api/utils/mysql_utils"
	"github.com/StefanPahlplatz/bookstore_utils-go/logger"
	"github.com/StefanPahlplatz/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, created_on, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, created_on, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, created_on, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, created_on, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Get the row.
	result := sqlStatement.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedOn, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Insert the user.
	insertResult, err := sqlStatement.Exec(user.FirstName, user.LastName, user.Email, user.CreatedOn, user.Status, user.Password)
	if err != nil {
		logger.Error("error when trying to insert user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}

	// Set the ID.
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get id after inserting user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Update the user
	_, updateErr := sqlStatement.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updateErr != nil {
		logger.Error("error when trying to execute update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Delete the user.
	if _, err := sqlStatement.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Get the rows.
	rows, err := sqlStatement.Query(status)
	if err != nil {
		logger.Error("error when trying to execute find by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedOn, &user.Status); err != nil {
			logger.Error("error when trying to execute find by status", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	// Prepare sql statement.
	sqlStatement, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer sqlStatement.Close()

	// Get the row.
	result := sqlStatement.QueryRow(user.Email, user.Password, StatusActive)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedOn, &user.Status); err != nil {
		if strings.Contains(err.Error(), mysql_utils.ErrorNoResults) {
			return rest_errors.NewNotFoundError("invalid credentials")
		}
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}

	return nil
}
