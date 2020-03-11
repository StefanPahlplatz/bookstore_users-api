package users_db

import (
	"database/sql"
	"fmt"
	"github.com/StefanPahlplatz/bookstore_users-api/logger"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_password = "mysql_users_password"
	mysql_users_host     = "mysql_users_host"
	mysql_users_schema   = "mysql_users_schema"
)

var (
	Client *sql.DB

	username = os.Getenv(mysql_users_username)
	password = os.Getenv(mysql_users_password)
	host     = os.Getenv(mysql_users_host)
	schema   = os.Getenv(mysql_users_schema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Error("unable to open mysql connection", err)
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		logger.Error("unable to ping mysql", err)
		panic(err)
	}
	log.Println("database sucessfully configured")
}
