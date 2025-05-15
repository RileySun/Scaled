package utils

import (
	"fmt"
	"log"
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

func NewDB(host, port, user, pass, dbName string) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)
	db, err := sql.Open("mysql", uri)
	
	if err != nil {
		log.Println(err, " - " + host + ":" + port + "/" + dbName + " - Connect")
		return nil, err
	}
	
	db.SetMaxIdleConns(100)//Connection pooling
	
	//Make sure connection is real
	err = db.Ping()
	if err != nil {
		log.Println(err, " - " + host + ":" + port + "/" + dbName + " - Ping")
		return nil, err
	}
	
	return db, nil
}

