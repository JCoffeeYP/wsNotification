package src

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ReadConfig() kafka.ConfigMap {
	// todo Использовать чтение из файла или из средозависимых переменных
	m := make(map[string]kafka.ConfigValue)
	m["bootstrap.servers"] = "localhost:9092"
	m["security.protocol"] = "PLAINTEXT"
	m["group.id"] = "TestGroup"
	m["auto.offset.reset"] = "earliest"

	return m
}

func ProcessMessages() {
	conf := ReadConfig()
	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		log.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "TestTopic" // todo Добавление средозависимой переменной
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("Failed to connect to topic: %s", err)
		os.Exit(1)
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
			c.Close()
		default:
			ev, err := c.ReadMessage(200 * time.Millisecond)
			if err != nil {
				continue
			}
			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			// todo добавить отправку данных в нужные Polls
		}
	}
}
