package mysql_utils

import (
	"errors"
	"github.com/StefanPahlplatz/bookstore_utils-go/rest_errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	ErrorNoResults = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoResults) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", err)
	}

	switch sqlErr.Number {
	case 1062:
		// Duplicate key
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error while processing request", errors.New("database error"))
}
