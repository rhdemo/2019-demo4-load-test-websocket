package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/sacOO7/gowebsocket"
	uuid "github.com/satori/go.uuid"
)

var playerID string

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	movement := os.Getenv("MOVEMENT")
	socket := gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error - ", err)
	}

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Received message - " + message)
		message = convertJSON(message)

		if message != "" {
			playerID = message
		}

	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()
	socket.SendBinary([]byte(`{ "type": "connection"}`))

	for {
		if playerID != "" {
			time.Sleep(5 * time.Second)
			motionPayload := createPayload(playerID, movement)
			socket.SendBinary(motionPayload)
		}
	}
}

func convertJSON(input string) string {
	var result map[string]interface{}
	json.Unmarshal([]byte(input), &result)
	playerId, _ := result["playerId"].(string)
	return playerId
}

func createPayload(playerID string, movement string) []byte {
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
		rand.Seed(time.Now().Unix())
		n = rand.Int() % len(moves)
	}

	moveSelected := "moves/" + moves[n]

	jsonFile, _ := os.Open(moveSelected)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)
	data["playerId"] = playerID
	data["uuid"] = uuid.Must(uuid.NewV4()).String()
	out, _ := json.Marshal(data)
	return out
}
