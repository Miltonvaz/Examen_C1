package servidorp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


func Run() {
	r := gin.Default()

	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.GET("/cambios", getCambios)
	r.GET("/send-to-replication", sendUserToReplication)

	srv := &http.Server{
		Addr:         ":4000",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  1 * time.Hour,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Error: Server Main hasn't begun")
	}
}
