package epochConsumer

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	exporter "epochConvertor/metrics"

	"github.com/Shopify/sarama"

	log "github.com/sirupsen/logrus"
)

const defaultGroup = "epochConsumer"
const defaultTopic = "input"

type ConsumerGroupMember struct {
	localMsgCh chan<- string
}

func (member *ConsumerGroupMember) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (member *ConsumerGroupMember) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (member *ConsumerGroupMember) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
		member.localMsgCh <- string(msg.Value)
		if _, ok := os.LookupEnv("CI"); !ok {
			incPrometuesCounter()
		}
	}
	return nil
}

func incPrometuesCounter() error {
	consumerGroup := getEnv("CONSUMER_GROUP", defaultGroup)
	topic := getEnv("CONSUMER_TOPIC", defaultTopic)
	metric := exporter.ConvertorMessageCounter.WithLabelValues(consumerGroup, topic)
	metric.Inc()
	return nil
}

func newConsumer() (sarama.ConsumerGroup, error) {
	consumerConfig := sarama.NewConfig()
	var brokers []string
	brokers = []string{"localhost:9092"}

	if value, ok := os.LookupEnv("CONSUMER_FETCH_MIN"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_FETCH_MIN: ", err)
		}
		consumerConfig.Consumer.Fetch.Min = int32(valuei)
	}
	if value, ok := os.LookupEnv("CONSUMER_FETCH_DEFAULT"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_FETCH_DEFAULT: ", err)
		}
		consumerConfig.Consumer.Fetch.Default = int32(valuei)
	}
	if value, ok := os.LookupEnv("CONSUMER_RETRY_BACKOFF"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_RETRY_BACKOFF: ", err)
		}
		consumerConfig.Consumer.Retry.Backoff = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	if value, ok := os.LookupEnv("CONSUMER_MAXWAITTIME"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_MAXWAITTIME: ", err)
		}
		consumerConfig.Consumer.MaxWaitTime = time.Duration(rand.Int31n(int32(valuei))) * time.Millisecond
	}
	if value, ok := os.LookupEnv("CONSUMER_MAXPROCESSINGTIME"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_MAXPROCESSINGTIME: ", err)
		}
		consumerConfig.Consumer.MaxProcessingTime = time.Duration(rand.Int31n(int32(valuei))) * time.Millisecond
	}
	consumerConfig.Consumer.Return.Errors = true
	if value, ok := os.LookupEnv("CONSUMER_RETURN_ERROR"); ok {
		if value == "false" {
			consumerConfig.Consumer.Return.Errors = false
		}
	}
	consumerConfig.Consumer.Offsets.AutoCommit.Enable = true
	if value, ok := os.LookupEnv("CONSUMER_OFFSETS_AUTO_COMMIT_ENABLE"); ok {
		if value == "false" {
			consumerConfig.Consumer.Offsets.AutoCommit.Enable = false
		}
	}
	if value, ok := os.LookupEnv("CONSUMER_OFFSETS_AUTO_COMMIT_INTERVAL"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_OFFSETS_AUTO_COMMIT_INTERVAL: ", value, err)
		}
		consumerConfig.Consumer.Offsets.AutoCommit.Interval = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	if value, ok := os.LookupEnv("CONSUMER_OFFSETS_INITIAL"); ok {
		if value == "oldest" {
			consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
		} else if value == "newest" {
			consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
		} else {
			log.Error("Bad! CONSUMER_OFFSETS_INITIAL: ", value)
		}
	}
	if value, ok := os.LookupEnv("CONSUMER_OFFSETS_RETRY_MAX"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_OFFSETS_RETRY_MAX: ", err)
		}
		consumerConfig.Consumer.Offsets.Retry.Max = valuei
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_SESSION_TIMEOUT"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_GROUP_SESSION_TIMEOUT: ", err)
		}
		consumerConfig.Consumer.Group.Session.Timeout = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_HEARTBEAT_INTERVAL"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_GROUP_HEARTBEAT_INTERVAL: ", err)
		}
		consumerConfig.Consumer.Group.Heartbeat.Interval = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_REBALANCE_STRATEGY"); ok {
		if value == "range" {
			consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
		} else if value == "sticky" {
			consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
		} else if value == "rr" {
			consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
		} else {
			log.Error("Bad! CONSUMER_GROUP_REBALANCE_STRATEGY: ", value)
		}
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_REBALANCE_TIMEOUT"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_GROUP_REBALANCE_TIMEOUT: ", err)
		}
		consumerConfig.Consumer.Group.Rebalance.Timeout = time.Duration(rand.Int31n(int32(valuei))) * time.Millisecond
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_REBALANCE_RETRY_MAX"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_GROUP_REBALANCE_RETRY_MAX: ", err)
		}
		consumerConfig.Consumer.Group.Rebalance.Retry.Max = valuei
	}
	if value, ok := os.LookupEnv("CONSUMER_GROUP_REBALANCE_RETRY_BACKOFF"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_GROUP_REBALANCE_RETRY_BACKOFF: ", err)
		}
		consumerConfig.Consumer.Group.Rebalance.Retry.Backoff = time.Duration(rand.Int31n(int32(valuei))) * time.Second
	}
	if value, ok := os.LookupEnv("CONSUMER_CLIENTID"); ok {
		consumerConfig.ClientID = value
	}
	if value, ok := os.LookupEnv("CONSUMER_CHANNELBUFFERSIZE"); ok {
		valuei, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Bad! CONSUMER_CHANNELBUFFERSIZE: ", err)
		}
		consumerConfig.ChannelBufferSize = valuei
	}
	if value, ok := os.LookupEnv("CONSUMER_VERSION"); ok {
		version, err := sarama.ParseKafkaVersion(value)
		if err != nil {
			log.Error("version incompatible: ", err)
		}
		consumerConfig.Version = version
	}
	if value, ok := os.LookupEnv("CONSUMER_BROKERS"); ok {
		brokers = []string{value}
	}
	group := defaultGroup
	if value, ok := os.LookupEnv("CONSUMER_GROUP"); ok {
		group = value
	}
	consumer, err := sarama.NewConsumerGroup(brokers, group, consumerConfig)
	return consumer, err
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
