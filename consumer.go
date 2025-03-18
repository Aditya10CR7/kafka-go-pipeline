package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	kafkaBroker = "54.211.215.139:9092" // Replace with your actual Kafka broker
	topic       = "task"
	groupID     = "consumer-group-1"

	awsRegion    = "us-east-1" // Replace with your AWS region
	s3BucketName = "task-aditya"
)

var awsAccessKey = os.Getenv("") // Set as environment variable
var awsSecretKey = os.Getenv("") // Set as environment variable

// uploadToS3 uploads JSON data to an S3 bucket
func uploadToS3(fileName string, data []byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("private"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	log.Printf("Uploaded %s to S3", fileName)
	return nil
}

// consumeMessages reads messages from Kafka and uploads them to S3
func consumeMessages() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup([]string{kafkaBroker}, groupID, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	ctx := context.Background()
	handler := ConsumerHandler{}

	for {
		err := consumer.Consume(ctx, []string{topic}, &handler)
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
			time.Sleep(2 * time.Second) // Retry on failure
		}
	}
}

// ConsumerHandler implements sarama.ConsumerGroupHandler
type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message: %s", string(msg.Value))

		// Generate a unique filename for each record
		fileName := fmt.Sprintf("stock_market_%d.json", time.Now().UnixNano())

		// Upload message to S3
		err := uploadToS3(fileName, msg.Value)
		if err != nil {
			log.Printf("S3 upload error: %v", err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	log.Println("Starting Kafka Consumer...")
	consumeMessages()
}
