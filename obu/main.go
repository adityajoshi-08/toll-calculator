package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"
	"toll-calculator/types"

	"github.com/gorilla/websocket"
)

const (
	sendInterval = time.Second * 5
)

const wsEndpoint = "ws://localhost:8080/ws"

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			// jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Error marshalling data: %v\n", err)
				continue
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)

	for i := 0; i < n; i++ {
		ids[i] = rand.IntN(math.MaxInt32)
	}
	return ids
}

func genLatLong() (float64, float64) {
	lat := generateCoordinates()
	long := generateCoordinates()
	return lat, long
}

func generateCoordinates() float64 {
	n := float64(rand.IntN(100) + 1)
	f := rand.Float64()
	return n + f
}
