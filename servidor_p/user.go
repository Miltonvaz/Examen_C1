package servidorp

import (
	"fmt"
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

type Cambio struct {
	Accion string `json:"accion"`
	User   User   `json:"user"`
}

var (
	bd      []User
	cambios []Cambio
)

func sendToReplicationServer(user User, accion string) {
	url := fmt.Sprintf("http://localhost:5000/replication?user_id=%d&name=%s&user=%s&accion=%s", user.ID, user.Name, user.User, accion)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error al enviar al servidor de replicación:", err)
		return
	}
	defer resp.Body.Close()
}

func sendUserToReplication(c *gin.Context) {
	if len(cambios) > 0 {
		lastChange := cambios[len(cambios)-1]
		sendToReplicationServer(lastChange.User, lastChange.Accion)
		c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario enviado a replicación", "usuario": lastChange.User})
	} else {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
	}
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser.ID = int64(len(bd) + 1)
	bd = append(bd, newUser)

	cambios = append(cambios, Cambio{Accion: "create", User: newUser})
	sendToReplicationServer(newUser, "create")

	c.JSON(http.StatusCreated, newUser)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userIndex := -1

	for i, user := range bd {
		if strconv.FormatInt(user.ID, 10) == id {
			userIndex = i
			cambios = append(cambios, Cambio{Accion: "delete", User: user})
			sendToReplicationServer(user, "delete")
			break
		}
	}

	if userIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"mensaje": "Usuario no encontrado"})
		return
	}

	bd = append(bd[:userIndex], bd[userIndex+1:]...)

	c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario eliminado"})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, user := range bd {
		if strconv.FormatInt(user.ID, 10) == id {
			updatedUser.ID = user.ID
			bd[i] = updatedUser

			cambios = append(cambios, Cambio{Accion: "update", User: updatedUser})
			sendToReplicationServer(updatedUser, "update")

			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"mensaje": "Usuario no encontrado"})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, bd)
}

func getCambios(c *gin.Context) {
	if len(cambios) == 0 {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
		return
	}

	response := cambios
	cambios = []Cambio{}

	c.JSON(http.StatusOK, response)
}
