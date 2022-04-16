package main

import (
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `gorm:"size:256"`
	Addr string `gorm:"size:256"`
	Age  int16  `gorm:"type:int"`
}

func OpenGORM() {
	if err := OpenMySQL(); err != nil {
		panic(err)
	}
	var conn, err = GetConn(context.TODO())
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: conn,
	}))

	if err != nil {
		panic(err)
	}

	if !db.Migrator().HasTable(&User{}) {
		db.Migrator().CreateTable(&User{})
	}

	var user = User{
		Name: "gan",
		Addr: "beijing",
		Age:  100,
	}
	var ret = db.Where(&User{
		Model: gorm.Model{},
		Name:  "gan",
	}).FirstOrCreate(&user, "name = ?", "gan")

	log.Printf("User:%+v, err:%+v", user, ret.Error)
}
