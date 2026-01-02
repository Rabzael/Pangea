package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Struct to represent a log entry
type LogEntry struct {
	Time       string `json:"time,omitempty"`
	RemoteAddr string `json:"remote_addr"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Status     string `json:"status"`
	ContentLen int64  `json:"content_len"`
	Error      string `json:"error"`
	Blocked    bool   `json:"blocked"`
}

// Global log file variable
var logFile *os.File

// Kafka producer
var kafkaProducer *kafka.Producer

func InitLogger() {
	// Open log file
	lf, err := os.OpenFile(Config.Log_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("ERROR: could not open log file: %s\n", err.Error())
		os.Exit(1)
	}
	logFile = lf

	// Initialize Kafka producer if enabled
	if Config.Kafka_enabled {
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": Config.Kafka_brokers})
		if err != nil {
			fmt.Printf("ERROR: could not create Kafka producer") // In case of failure, continue without Kafka
			kafkaProducer = nil
		}
		kafkaProducer = p
	} else {
		kafkaProducer = nil
		fmt.Printf("INFO: Kafka logging is disabled.\n")
	}
}

// Close the logger
func CloseLogger() {
	if logFile != nil {
		logFile.Sync()
		logFile.Close()
	}

	if kafkaProducer != nil {
		kafkaProducer.Flush(5000)
		kafkaProducer.Close()
	}
}

func logAppend(entry LogEntry) error {
	ljs, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("ERROR: could not marshal log entry: %s\n", err.Error())
		return err
	}

	err = logFileWrite(ljs)
	if err != nil {
		return err
	}

	if kafkaProducer != nil {
		err = logKafkaSend(ljs)
		if err != nil {
			return err
		}
	}

	return nil
}

// Append a log entry to the log file
func logFileWrite(entry []byte) error {
	_, err := logFile.WriteString(string(entry) + "\n")
	if err != nil {
		fmt.Printf("ERROR: could not write logs file: %s\n", err.Error())
		return err
	}
	return nil
}

// Send log entry to Kafka
func logKafkaSend(entry []byte) error {
	kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Config.Kafka_topic, Partition: kafka.PartitionAny},
		Value:          []byte(entry),
	}, nil)
	return nil
}

/** Generate logs for different scenarios **/
func logBase(req *http.Request) LogEntry {
	return LogEntry{
		Time:       time.Now().Format(time.RFC3339),
		RemoteAddr: req.RemoteAddr,
		Method:     req.Method,
		URL:        req.URL.String(),
	}
}

func LogOk(req *http.Request, res *http.Response) {
	l := logBase(req)
	l.Status = res.Status
	l.ContentLen = res.ContentLength

	logAppend(l)
}

func LogError(req *http.Request, err error) {
	l := logBase(req)
	l.Error = err.Error()

	logAppend(l)
}

func LogBlocked(req *http.Request) {
	l := logBase(req)
	l.Blocked = true

	logAppend(l)
}
