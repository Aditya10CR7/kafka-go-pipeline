package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
)

const (
	kafkaBroker = "54.211.215.139:9092"
	topic       = "task"
	apiURL      = "https://files.polygon.io/api/v1/"
)

// FetchDataFromAPI fetches the data from the provided API
func FetchDataFromAPI(url string) ([]map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received non-OK HTTP status %v", resp.StatusCode)
	}

	var data []map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding API response: %v", err)
	}

	return data, nil
}

func main() {
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	// Fetch data from the API
	data, err := FetchDataFromAPI(apiURL)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	for {
		randomIndex := rand.Intn(len(data))
		message, _ := json.Marshal(data[randomIndex])

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}

		_, _, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			fmt.Println("Message sent:", string(message))
		}

		time.Sleep(2 * time.Second) // Simulate real-time streaming
	}
}
