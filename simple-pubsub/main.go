package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type contents struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

var ps = NewPubSub()

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, should be POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var c contents
	if err := json.Unmarshal(body, &c); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if c.Topic == "" || c.Message == "" {
		http.Error(w, "Topic and message are required", http.StatusBadRequest)
		return
	}
	ps.Publish(c.Topic, c.Message)
	log.Printf("Successfully published message to topic: %v", c.Topic)
	w.WriteHeader(http.StatusOK)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, should be GET", http.StatusMethodNotAllowed)
		return
	}

	// Check the topi from query parameters
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	// Clean up when done
	ch := make(chan string, 1)
	defer func() {
		ps.mu.Lock()
		defer ps.mu.Unlock()
		for _, subscribers := range ps.topics {
			for i, c := range subscribers {
				if c == ch {
					ps.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
					close(c)
				}
			}
		}
	}()

	// Listens for messages and writes them to the response
	wait := make(chan struct{})
	go func() {
		defer func() { wait <- struct{}{} }()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				_, err := w.Write([]byte("\n" + msg + "\n"))
				if err != nil {
					return
				}
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}
		}
	}()
	_, err := w.Write([]byte(fmt.Sprintf("Subscribed to topic: %v, waiting for messages...", topic)))
	if err != nil {
		log.Println(err)
		return
	}
	w.(http.Flusher).Flush()
	// Subscribe to the topic, adds the channel to the list of subscribers
	ps.Subscribe(topic, ch)

	<-wait
}

type PubSub struct {
	topics map[string][]chan string // key: topic, val: list of subscribers
	mu     sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		topics: make(map[string][]chan string),
	}
}

func (ps *PubSub) Publish(topic string, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if subscribers, ok := ps.topics[topic]; ok {
		for _, ch := range subscribers {
			ch <- message
		}
	}
}

func (ps *PubSub) Subscribe(topic string, ch chan string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.topics[topic] = append(ps.topics[topic], ch)
}

func main() {
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

	log.Println("Starting pubsub server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
