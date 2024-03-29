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
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/go_test?parseTime=true")

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
	defer cancel()
	// var conn, err = GetConn(todo)
	// if err != nil {
	// 	log.Panicf("Open Conn error:%v", err)
	// }
	// defer func(conn *sql.Conn) {
	// 	err := conn.Close()
	// 	if err != nil {
	// 		log.Printf("Close error:%v", err)
	// 	}
	// }(conn)
	var rows, err = DB.QueryContext(todo, "select now()")
	if err != nil {
		log.Printf("query error:%v", err)
		return
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Printf("close query row error %v", closeErr)
		}
	}()
	for rows.Next() {
		var val string
		_ = rows.Scan(&val)
		log.Println(val)
	}
}
