package notification_svc

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type NotificationService struct {
	reader KafkaReader
}

func NewNotificationService(reader KafkaReader) *NotificationService {
	return &NotificationService{reader: reader}
}

func (n *NotificationService) ProcessMessages(ctx context.Context) {
	for {
		msg, err := n.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		log.Printf("Received message: %s", string(msg.Value))

		// Имитация отправки уведомления
		err = SendNotification(string(msg.Value))
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}
}

func (n *NotificationService) Close() {
	n.reader.Close()
}

func SendNotification(message string) error {
	log.Printf("Notification sent: %s", message)
	return nil
}

func main() {
	broker := "localhost:9092"
	topic := "booking-events"
	groupID := "notification-group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	notificationSvc := NewNotificationService(reader)
	defer notificationSvc.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go notificationSvc.ProcessMessages(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Notification Service shutting down")
}
