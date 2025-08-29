package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"toll-calculator/aggregator/client"
	"toll-calculator/types"

	"google.golang.org/grpc"
)

// Transporter lauda lehsun
func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "listen address of the GRPC server")
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpcListenAddr, svc)
	time.Sleep(time.Second * 5);
	c, err := client.NewGRPCClient(*grpcListenAddr);
	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.AggregatorClient.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 1,
		Value: 1.1,
		Unix: time.Now().Unix(),
	}); err != nil {
		log.Fatal(err)
	}
	makeHTTPTransport(*httpListenAddr, svc)
	fmt.Println("Working fyne af")
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("Starting HTTP server on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	// make a tcp listener
	fmt.Println("Starting GRPC server on", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu id"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
