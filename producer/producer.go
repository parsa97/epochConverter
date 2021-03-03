package rfcProducer

import (
	"context"
	exporter "epochConvertor/metrics"
	"os"
	"strconv"
	"time"

	cmap "github.com/orcaman/concurrent-map"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func StartProducer(ctx context.Context, errCH chan error, msgCH chan string) error {
	log.Info("Server Start Producing")
	topic := "output"
	if value, ok := os.LookupEnv("PRODUCER_TOPIC"); ok {
		topic = value
	}
	producer, err := newProducer()
	if err != nil {
		log.Error("Could not create producer: ", err)
		errCH <- err
		return err
	}
	var partitionProduced = cmap.New()
	go func() {
		for {
			time.Sleep(time.Millisecond)
			epochMsg, _ := <-msgCH
			msg := prepareMessage(topic, epochMsg)
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Error(err)
				time.Sleep(1 * time.Second)
			}
			counter, _ := partitionProduced.Get(string(partition))
			if counter == nil {
				partitionProduced.Set(string(partition), 1)
			} else {
				counted := counter.(int)
				counted++
				partitionProduced.Set(string(partition), counted)
				log.Trace("Counter: ", " Message: ", epochMsg, " topic: ", topic, " partition: ", partition, " offset: ", offset)
				metric := exporter.ProducedMessageCounter.WithLabelValues(strconv.Itoa(int(partition)), topic)
				metric.Inc()
			}
		}
	}()
	<-ctx.Done()
	log.Info("Producer Stopped")

	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancel()
	}()
	return nil
}

func prepareMessage(topic, epochMsg string) *sarama.ProducerMessage {
	msgi, err := strconv.Atoi(epochMsg)
	if err != nil {
		log.Error(err)
	}
	message := time.Unix(0, int64(msgi)*int64(time.Millisecond)).Format(time.RFC3339Nano)
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}
	return msg
}
