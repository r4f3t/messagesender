package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/r4f3t/messagesender/helper"
	"gorm.io/gorm"
)

type MessageService interface {
	StartProcessingMessages() error
	StopProcessingMessages() error
	GetSentMessages() ([]Message, error)
}

type service struct {
	messageRepository MessageRepository
	redisDb           *redis.Client
}

func NewService(messageRepository MessageRepository, redisDb *redis.Client) MessageService {
	return &service{
		messageRepository: messageRepository,
		redisDb:           redisDb,
	}
}

var (
	processRunning = false
	processStop    = make(chan bool)
	processMutex   sync.Mutex
)

// StartProcessingMessages starts the periodic message processing
func (receiver *service) StartProcessingMessages() error {
	processMutex.Lock()
	defer processMutex.Unlock()

	if processRunning {
		return errors.New("Message processing is already running.")
	}

	processRunning = true
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		log.Println("Message processing started.")
		for {
			select {
			case <-processStop:
				log.Println("Message processing stopped.")
				return
			case <-ticker.C:
				log.Println("Fetching unsent messages...")
				limit := 2
				messages, err := receiver.messageRepository.FetchMessages(&limit, false)
				if err != nil {
					log.Printf("Error fetching messages: %v", err)
					continue
				}

				for _, message := range messages {

					err := receiver.messageRepository.GetDBInstance().Transaction(func(tx *gorm.DB) error {

						if errInternal := sendMessage(message); err != nil {
							log.Printf("Error sending message ID %d: %v", message.ID, err)
							return errInternal
						}

						if errInternal := receiver.messageRepository.MarkMessageAsSent(message.ID, tx); err != nil {
							log.Printf("Error marking message ID %d as sent: %v", message.ID, err)
							return errInternal
						}

						if errInternal := helper.CacheMessage(message.ID, receiver.redisDb); err != nil {
							log.Printf("Error caching message ID %d: %v", message.ID, err)
							return errInternal
						}

						return nil
					})

					if err != nil {
						log.Printf("Transaction failed for message ID %d: %v", message.ID, err)
						continue
					}
				}
			}
		}
	}()

	return nil
}

// StopProcessingMessages stops the periodic message processing
func (receiver *service) StopProcessingMessages() error {
	processMutex.Lock()
	defer processMutex.Unlock()

	if !processRunning {
		return errors.New("Message processing is not running.")
	}

	processRunning = false
	processStop <- true

	return nil
}

// Get Sent messages retrive all sent messages
func (receiver *service) GetSentMessages() ([]Message, error) {
	return receiver.messageRepository.FetchMessages(nil, true)
}

// SendMessage sends a message to the specified webhook
func sendMessage(message Message) error {
	url := os.Getenv("WEBHOOK_URL")
	payload := map[string]string{
		"to":      message.To,
		"content": message.Content,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
