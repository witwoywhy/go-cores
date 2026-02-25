package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	MessageChannel chan *Message
	SetupChannel   chan struct{}
}

func NewConsumerGroupHandler(
	messageChannel chan *Message,
	setupChannel chan struct{},
) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		MessageChannel: messageChannel,
		SetupChannel:   setupChannel,
	}
}

func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group session setup")
	h.SetupChannel <- struct{}{}
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group session cleanup")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.MessageChannel <- &Message{
			Key:       string(msg.Key),
			Value:     string(msg.Value),
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
