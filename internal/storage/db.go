package storage

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Unparalleled-Calvin/gserver/internal/schema"
	"github.com/Unparalleled-Calvin/gserver/internal/settings"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db        *sql.DB
	mysqlOnce sync.Once
)

func GetDB() *sql.DB {
	mysqlOnce.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", settings.MySQLUser, settings.MySQLPassword, settings.MySQLAddr, settings.MySQLDB)
		var err error
		if db, err = sql.Open("mysql", dsn); err != nil {
			panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
		}
	})
	return db
}

func InsertUser(userRegister schema.UserRegister) error {
	sql := "INSERT INTO user (id, username, password_hash) VALUES (?, ?, ?)"
	_, err := GetDB().Exec(sql, userRegister.ID, userRegister.UserName, userRegister.Password)
	return err
}
