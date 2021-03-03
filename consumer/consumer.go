package epochConsumer

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func StartConsumer(ctx context.Context) {
	consumerTopic := "input"
	if value, ok := os.LookupEnv("CONSUMER_TOPIC"); ok {
		consumerTopic = value
	}
	consumer, err := newConsumer()
	if err != nil {
		ConsumerErrCH <- err
		return
	}

	go func() {
		for err := range consumer.Errors() {
			ConsumerErrCH <- err
			return
		}
	}()

	if err := getEpochMessage(ctx, consumer, consumerTopic); err != nil {
		log.Error("failed to Consume: ", err)
		ConsumerErrCH <- err
	}
}

func getEpochMessage(ctx context.Context, consumer sarama.ConsumerGroup, topic string) error {
	log.Info("Server Start Consuming")
	//var partitionConsumed = cmap.New()
	member := ConsumerGroupMember{MsgCh: EpochTimes}
	go func() {
		for {
			err := consumer.Consume(ctx, strings.Split(topic, ","), &member)
			if err != nil {
				//log.Println(err)
				ConsumerErrCH <- err
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-ctx.Done()
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	return nil
}
