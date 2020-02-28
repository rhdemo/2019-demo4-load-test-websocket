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

	// uuid "github.com/satori/go.uuid"
)

type user struct {
	playerID string
	gameID   string
	score  float64
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

	socket := gowebsocket.New(socketAddress)

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
		if u.socket.IsConnected == true {
			u.socket.SendBinary([]byte(`{ "type": "init", "bot": "true"}`))
		}
	}

	u.socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		messageResult := convertJSON(message)

		if messageResult["type"].(string) == "player-configuration" {
			player := messageResult["player"].(map[string]interface{})
			u.playerID = player["id"].(string)
			u.gameID = player["gameId"].(string)
			u.score = player["score"].(float64)
			log.Println(strconv.Itoa(clientNumber) + " became " + u.playerID)
			log.Println( u.playerID + " score: " + fmt.Sprintf("%f", u.score)
		}
	}
	u.socket.Connect()

	for {
		time.Sleep(5 * time.Second)
		if u.playerID != "" && u.socket.IsConnected == true {
			selectedMovement, movementLog := selectRandomMovement(movement)
			motionPayload := createPayload(u.playerID, selectedMovement)
			log.Println(u.playerID + " is sending a " + movementLog)
			u.socket.SendBinary(motionPayload)
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

		reconnectPayload := fmt.Sprintf(`{ "type": "init", "bot": "true", "gameId": "%s", "playerId": "%s" }`, u.playerID, u.gameID)
		u.socket = gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))
		u.socket.Connect()
		u.socket.SendBinary([]byte(reconnectPayload))
		break
	}
}

func selectRandomMovement(movement string) (string, string) {
	var n int
	moves := []string{
		"bad.json",
		"good.json",
	}
	switch movement {
	case "BAD":
		n = 0
	case "GOOD":
		n = 1
	// case "NO-BAD-DANCE":
	// 	n = rand.Int() % (len(moves)-1)
	case "RANDOM":
		n = rand.Int() % len(moves)
	}

	moveSelected := "guesses/" + moves[n]
	movelog := strings.TrimSuffix(moves[n], filepath.Ext(moves[n]))
	return moveSelected, movelog
}

func createPayload(playerID string, movement string) []byte {


	jsonFile, _ := os.Open(movement)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)
	data["playerId"] = playerID
	out, _ := json.Marshal(data)
	return out
}
