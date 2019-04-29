package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sacOO7/gowebsocket"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	playerID string
	gameID   string
	socket   gowebsocket.Socket
}

var exit = make(chan bool)

func main() {
	users, _ := strconv.Atoi(os.Getenv("USERS"))
	movement := os.Getenv("MOVEMENT")
	socketAddress := os.Getenv("SOCKET_ADDRESS")

	for clientNumber := 1; clientNumber <= users; clientNumber++ {
		go handleSocket(clientNumber, socketAddress, movement)
	}
	<-exit // This blocks until the exit channel receives some input
	fmt.Println("Done.")

}

func handleSocket(clientNumber int, socketAddress string, movement string) {

	socket := gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))

	u := user{playerID: "", gameID: "", socket: socket}

	u.socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println(u.playerID+" received connect error - ", err)
	}

	u.socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println(u.playerID+" disconnected - ", err)
		reconnect(u)
	}

	u.socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println(strconv.Itoa(clientNumber) + " connected")
		u.socket.SendBinary([]byte(`{ "type": "connection"}`))
	}

	u.socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		messageResult := convertJSON(message)

		if messageResult["type"].(string) == "configuration" {
			u.playerID = messageResult["playerId"].(string)
			u.gameID = messageResult["gameId"].(string)
			log.Println(strconv.Itoa(clientNumber) + " became " + u.playerID)
		} else if messageResult["type"].(string) == "motion_feedback" {
			log.Println(u.playerID + " received message - " + message)
		}
	}
	u.socket.Connect()

	for {
		time.Sleep(5 * time.Second)
		if u.playerID != "" && u.socket.IsConnected == true {
			selectedMovement, movementLog := selectRandomMovement(movement)
			motionPayload := createPayload(u.playerID, selectedMovement)
			log.Println(u.playerID + "is about to send a " + movementLog)
			u.socket.SendBinary(motionPayload)
		} else if u.socket.IsConnected == false {
			reconnect(u)
		}
	}
	// This is will not happen because of time infinite loop (#TODO change to duration)
	exit <- true
}

func convertJSON(input string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(input), &result)
	return result
}

func reconnect(u user) {
	for {
		log.Println("User: " + u.playerID + " is reconnecting")
		time.Sleep(2 * time.Second)
		reconnectPayload := `{ "type": "connection"` + `gameId:` + u.gameID + `playerId:` + u.playerID + `}`
		u.socket = gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))
		u.socket.SendBinary([]byte(reconnectPayload))
		break
	}
}

func selectRandomMovement(movement string) (string, string) {
	var n int
	moves := []string{
		"floss.json",
		"fever.json",
		"roll.json",
		"shake.json",
		"x.json",
		"circle.json",
		"bad-move.json",
	}
	switch movement {
	case "FLOSS":
		n = 0
	case "FEVER":
		n = 1
	case "ROLL":
		n = 2
	case "SHAKE":
		n = 3
	case "X":
		n = 4
	case "CIRCLE":
		n = 5
	case "BAD":
		n = 6
	case "RANDOM":
		n = rand.Int() % len(moves)
	}

	moveSelected := "moves/" + moves[n]
	movelog := strings.TrimSuffix(moves[n], filepath.Ext(moves[n]))
	return moveSelected, movelog
}

func createPayload(playerID string, movement string) []byte {
	jsonFile, _ := os.Open(movement)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)
	data["playerId"] = playerID
	data["uuid"] = uuid.Must(uuid.NewV4()).String()
	out, _ := json.Marshal(data)
	return out
}
