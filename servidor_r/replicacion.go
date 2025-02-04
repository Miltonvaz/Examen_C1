package servidorr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	User string `json:"user"`
}


var bdReplication []User

func getReplicatedUsers(c *gin.Context) {
	c.JSON(http.StatusOK, bdReplication)
}


func replicateUser(c *gin.Context) {
	userIDStr := c.Query("user_id")
	name := c.Query("name")
	user := c.Query("user")


	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	newUser := User{ID: userID, Name: name, User: user}
	bdReplication = append(bdReplication, newUser)

	c.JSON(http.StatusOK, gin.H{"message": "User replicated successfully"})
}

func Do_before_init() {
	for {
		
		response, err := http.Get("http://localhost:4000/users") 
		if err != nil {
			fmt.Println("Error fetching replication data:", err)
			time.Sleep(5 * time.Second) 
			continue
		}

		if response.StatusCode == http.StatusOK {
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			var updatedUsers []User
			if err := json.Unmarshal(body, &updatedUsers); err != nil {
				fmt.Println("Error parsing replication data:", err)
				continue
			}

			bdReplication = updatedUsers
			fmt.Println("Replication data updated!")
		}

		time.Sleep(10 * time.Second) 
	}
}
