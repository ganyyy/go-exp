package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func OpenMySQL() error {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/go_test?parseTime=true")

	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	DB = db
	return nil
}

func GetConn(done context.Context) (*sql.Conn, error) {
	return DB.Conn(done)
}

func TestMySql() {

	if err := OpenMySQL(); err != nil {
		panic(err)
	}

	var todo, cancel = context.WithCancel(context.TODO())
	var conn, err = GetConn(todo)
	if err != nil {
		log.Panicf("Open Conn error:%v", err)
	}
	defer cancel()
	defer func(conn *sql.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Close error:%v", err)
		}
	}(conn)
	var rows, _ = conn.QueryContext(todo, "select now()")
	for rows.Next() {
		var val string
		_ = rows.Scan(&val)
		log.Println(val)
	}
}
