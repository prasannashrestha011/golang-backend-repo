package method

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"projec/hashing"
	"projec/jwttoken"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// connecting to database table
var db, err = sql.Open("mysql", "root:9843@tcp(localhost:3306)/testdb")

// data structure
type UserKey struct {
	Id       *int   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserList struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Insert_key(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	var user_data UserKey
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, &user_data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	//inserting user key to database
	insertQuery := "INSERT INTO data key (username,password) VALUES(?,?)"
	hashedPassword, err := hashing.HashPassword(user_data.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, err = db.Exec(insertQuery, user_data.Username, hashedPassword)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	fmt.Println("data insertion sucessful")
	c.String(http.StatusOK, "data insertion sucessful")
}
func Get_All_Key(c *gin.Context) {
	var userList []UserKey
	//fetching all data
	fetchQuery := "SELECT * FROM  datakey"
	rows, err := db.Query(fetchQuery)

	if err != nil {
		fmt.Println("data fetching failed..")
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	for rows.Next() {
		var user_data UserKey
		err := rows.Scan(&user_data.Id, &user_data.Username, &user_data.Password)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		userList = append(userList, user_data)
	}

	c.JSON(http.StatusOK, userList)
}

func Authenticate_handler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	var user_data UserKey
	var fetch_user_data UserKey
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, &user_data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//fetching the user from database
	fetchQuery := "SELECT * FROM datakey WHERE username=?"

	rows, err := db.Query(fetchQuery, user_data.Username)
	if err != nil {
		fmt.Println("failed to fetch the data from db")
		return
	}
	for rows.Next() {
		err := rows.Scan(&fetch_user_data.Id, &fetch_user_data.Username, &fetch_user_data.Password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	isUserAuthenticated := hashing.ComparePassword(fetch_user_data.Password, user_data.Password)
	if isUserAuthenticated {
		tokenString, err := jwttoken.CreateToken(fetch_user_data.Username)

		if err != nil {

			fmt.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
		})
	} else {
		c.String(http.StatusUnauthorized, "User is not authenticated")
	}
}
