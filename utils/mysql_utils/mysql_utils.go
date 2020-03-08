package mysql_utils

import (
	"github.com/StefanPahlplatz/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	errorNoResults = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoResults) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		// Duplicate key
		return errors.NewBadRequestError("invalid data")
	}
	return errors.NewInternalServerError("error while processing request")
}
