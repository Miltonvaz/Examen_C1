package servidorr

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	r.GET("/replication", replicateUser) 

	if err := r.Run(":5000"); err != nil {
		fmt.Println("Error: Replication Server hasn't begun")
	}
}
