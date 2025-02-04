package servidorp

import (
	"fmt"
	"net/http"
	"net/url"
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
	url := fmt.Sprintf("http://localhost:5000/replication?user_id=%d&name=%s&user=%s&accion=%s",
		user.ID,
		url.QueryEscape(user.Name),
		url.QueryEscape(user.User),
		url.QueryEscape(accion),
	)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error al enviar al servidor de replicación:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error replicando el usuario:", resp.Status)
	}
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
