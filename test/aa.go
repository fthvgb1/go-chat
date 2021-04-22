package main

import "github.com/gin-gonic/gin"

func main() {

	r := gin.Default()
	r.GET("/aa", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong林蝇是",
		})
	})
	r.Run("0.0.0.0:9888")

}
