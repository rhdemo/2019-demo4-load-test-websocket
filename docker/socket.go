package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sacOO7/gowebsocket"
	uuid "github.com/satori/go.uuid"
)

var motionPayload []byte

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error - ", err)
	}

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Received message - " + message)
		playerID := convertJSON(message)

		if playerID != "" {
			motionPayload = createPayload(playerID)
		}

	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()
	socket.SendBinary([]byte(`{ "type": "connection"}`))

	//TODO Add Timer
	//TODO Separate on Go routines for user
	//TODO Bad Dances
	//

	for {
		if motionPayload != nil {
			time.Sleep(5 * time.Second)
			socket.SendBinary(motionPayload)
		}
		// select {
		// case <-interrupt:
		// 	log.Println("interrupt")
		// 	socket.Close()
		// 	return
		// }
	}
}

func convertJSON(input string) string {
	var result map[string]interface{}
	json.Unmarshal([]byte(input), &result)
	playerId, _ := result["playerId"].(string)
	return playerId
}

func createPayload(playerID string) []byte {
	// TODO add other types here
	jsonFile, _ := os.Open("floss.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)
	data["playerId"] = playerID
	data["uuid"] = uuid.Must(uuid.NewV4()).String()
	out, _ := json.Marshal(data)
	return out
}
