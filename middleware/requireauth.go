package middleware

import (
	"fmt"
	"net/http"
	"projec/jwttoken"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	fmt.Println(tokenString)
	if err != nil {
		fmt.Println(err.Error(), "is here")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "cookie not present,please login in",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = jwttoken.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err.Error(), "is from re")

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}
