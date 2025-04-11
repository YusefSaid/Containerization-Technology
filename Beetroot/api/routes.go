package api

import (
	"os"
	"fmt"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time":time.Now(),
	})
}

func GetEnv(c *gin.Context) {
	tenv := "Cyb3r53cur1ty"
	c.JSON(http.StatusOK, gin.H{
		tenv: os.Getenv(tenv),
	})
}

func GetFile(c *gin.Context) {
	fn := ".secret"
	cont, err := os.ReadFile("/dev/shm/"+fn)
	if err != nil{
		panic("missing "+fn)
	}
	c.JSON(http.StatusOK, gin.H{
		fn: string(cont),
	})
}

func root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Online",
	})
}


func StartAPI(Port int) {
	router := gin.Default()
    	router.GET("/", root)
	v1 := router.Group("/api/v1")
	{
		time := v1.Group("/time")
		{
			time.GET("/", GetTime)
		}
		env := v1.Group("/env")
		{
			env.GET("/", GetEnv)
		}
		File := v1.Group("/file")
		{
			File.GET("/", GetFile)
		}
	}
	router.Run(fmt.Sprintf(":%d", Port))
}
