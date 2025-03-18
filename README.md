# **Kafka Producer and Consumer - README**

## **Overview**
This project implements a real-time data streaming system using Kafka. It consists of:
1. A **Kafka Producer** that fetches data from an external API and sends it to a Kafka topic.
2. A **Kafka Consumer** that reads messages from the Kafka topic and uploads them to an AWS S3 bucket.

---

## **Setup and Prerequisites**

### **1. Install Dependencies**
Ensure you have the following installed:
- Go (Golang)
- Kafka
- AWS CLI (for S3 access)
- Required Go modules:
  ```sh
  go get github.com/Shopify/sarama
  go get github.com/aws/aws-sdk-go
  ```

### **2. Kafka Configuration**
Ensure Kafka is running and configured with the following:
- **Kafka Broker**: `54.211.215.139:9092`
- **Kafka Topic**: `task`

### **3. AWS S3 Configuration**
Set up an S3 bucket for storing consumed messages:
- **S3 Bucket Name**: `task-aditya`
- **AWS Region**: `us-east-1`

Set environment variables for AWS credentials:
```sh
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
```

---

## **Kafka Producer**
### **Functionality**
- Fetches stock market data from an external API.
- Randomly selects an entry and sends it to Kafka.
- Runs in an infinite loop, sending data every 2 seconds.

### **Run Producer**
```sh
go run producer.go
```

---

## **Kafka Consumer**
### **Functionality**
- Reads messages from the Kafka topic `task`.
- Generates a unique JSON file for each message.
- Uploads the file to an AWS S3 bucket.

### **Run Consumer**
```sh
go run consumer.go
```

---

## **Expected Output**
### **Producer Output**
```
Message sent: {"symbol":"AAPL", "price":175.23}
```

### **Consumer Output**
```
Received message: {"symbol":"AAPL", "price":175.23}
Uploaded stock_market_1710963847.json to S3
```

---

## **Conclusion**
This project successfully demonstrates real-time data streaming using Kafka and AWS S3. The producer fetches and streams data, while the consumer processes and stores it securely.

