package main

import (
	"database/sql"
	"fmt"
	"log"
	"projec/method"
	"projec/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db, err = sql.Open("mysql", "root:9843@tcp(localhost:3306)/testdb")

func main() {

	r := gin.Default()

	r.POST("/register", method.Insert_key)
	r.GET("/get-all-key", middleware.RequireAuth, method.Get_All_Key)
	r.POST("/login", method.Authenticate_handler)

	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error:", err)
	} else {
		fmt.Println("connection sucessfull")
	}
	r.Run(":8080")
}
