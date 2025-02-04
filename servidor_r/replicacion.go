package servidorr

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	User string `json:"user"`
}

var bdReplication []User

func getReplicatedUsers(c *gin.Context) {
	userIDStr := c.DefaultQuery("user_id", "")
	name := c.DefaultQuery("name", "")
	user := c.DefaultQuery("user", "")
	accion := c.DefaultQuery("accion", "")

	if userIDStr != "" && name != "" && user != "" && accion != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		newUser := User{ID: userID, Name: name, User: user}
		bdReplication = append(bdReplication, newUser)
		fmt.Println("User replicated:", newUser)
	}

	c.JSON(http.StatusOK, bdReplication)
}

