package main

import (
	"context"
	consumer "epochConvertor/consumer"
	exporter "epochConvertor/metrics"
	producer "epochConvertor/producer"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(logLevel())
}

var errCH chan error

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go exporter.Exporter()
	go consumer.StartConsumer(ctx)
	go producer.StartProducer(ctx)
	go func() {
		select {
		case err := <-exporter.ExporterErrCH:
			log.Println(err)
			c <- os.Interrupt
		case err := <-consumer.ConsumerErrCH:
			log.Println(err)
			c <- os.Interrupt
		case err := <-producer.ProducerErrCH:
			log.Println(err)
			c <- os.Interrupt
		}
	}()
	oscall := <-c
	log.Debug("system call: ", oscall)
	cancel()
}
