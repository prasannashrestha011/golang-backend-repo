package filemanagement

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	fmt.Println(fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to read the file",
		})
		return
	}
	//generating id for the files
	randomNumberString := strconv.Itoa(rand.Intn(9999))
	filePath := randomNumberString + "-" + fileHeader.Filename

	defer file.Close()
	err = c.SaveUploadedFile(fileHeader, "./public/"+filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error saving the file",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "file saved to the path",
	})
}
