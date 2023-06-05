package main

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/dhall0803/go_url_shortener/backend/lib/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func randomStringGenerator(length int) string {
	count := 0
	randomString := ""
	for count < length {
		// Generate a random number between 0 and 25
		randomNumber := rand.Intn(25) + 60

		// Convert the number to a letter
		randomLetter := string(rune(randomNumber))

		// Append the letter to the string
		randomString += randomLetter
		count++
	}

	return randomString
}

func createShortUrl() string {
	return os.Getenv("URL") + randomStringGenerator(6)
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a simple route
	router.GET("/url", func(c *gin.Context) {
		userId := c.Query("userid")
		longUrl := c.Query("longurl")

		if userId == "" || longUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Missing required parameters",
				"data":    nil,
			})
		}

		result, err := database.GetShortUrl(userId, longUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"data":    nil,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    result,
		})
	})

	router.POST("/url", func(c *gin.Context) {
		var newShortUrl database.ShortUrl
		err := c.BindJSON(&newShortUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"data":    nil,
			})
		}

		newShortUrl.ShortUrl = createShortUrl()
		newShortUrl.Id = uuid.New().String()

		err = database.CreateShortUrl(newShortUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"data":    nil,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    newShortUrl,
		})
	})

	// Start the server
	router.Run(":8080")
}
