package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"toll-calculator/types"
)

type DataReceiver struct {
	msgChan chan types.OBUData
	conn    *websocket.Conn
	prod    DataProducer
}

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal("Error creating data receiver:", err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	fmt.Println("Starting Data Receiver on :8080")
	http.ListenAndServe(":8080", nil)
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic string = "obu_data"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)

	return &DataReceiver{
		msgChan: make(chan types.OBUData, 128),
		prod:    p,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU Client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Error reading json:", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("Error producing data to kafka:", err)
		}
	}
}
