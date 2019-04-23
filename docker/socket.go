package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/sacOO7/gowebsocket"
	uuid "github.com/satori/go.uuid"
)

func main() {

	var wg sync.WaitGroup
	users, _ := strconv.Atoi(os.Getenv("USERS"))

	wg.Add(users)

	for clientNumber := 0; clientNumber < users; clientNumber++ {
		go func(i int) {
			defer wg.Done()
			var playerID string

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)

			movement := os.Getenv("MOVEMENT")
			socket := gowebsocket.New(os.Getenv("SOCKET_ADDRESS"))

			socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
				log.Fatal("Received connect error - ", err)
			}

			socket.OnConnected = func(socket gowebsocket.Socket) {
				log.Println("Connected to server")
				socket.SendBinary([]byte(`{ "type": "connection"}`))
				for {
					if playerID != "" {
						time.Sleep(5 * time.Second)
						motionPayload := createPayload(playerID, movement)
						socket.SendBinary(motionPayload)
					}
				}
			}

			socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
				log.Println("Received message - " + message)
				messageResult := convertJSON(message)

				if messageResult["type"].(string) == "configuration" && messageResult["gameState"].(string) == "active" {
					playerID = messageResult["playerId"].(string)
				}

			}

			socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
				log.Println("Disconnected from server ")
				return
			}

			socket.Connect()

		}(clientNumber)
	}

	wg.Wait()

}

func convertJSON(input string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(input), &result)
	return result
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
