package epochConsumer

import (
	"os"
	"testing"
)

func TestNewProducerBad(t *testing.T) {
	_, err := newConsumer()
	if err != nil {
		t.Error()
	}
}

func TestNewConsumer(t *testing.T) {
	os.Setenv("CONSUMER_FETCH_MIN", "1000")
	os.Setenv("CONSUMER_FETCH_DEFAULT", "10")
	os.Setenv("CONSUMER_RETRY_BACKOFF", "10")
	os.Setenv("CONSUMER_MAXWAITTIME", "10")
	os.Setenv("CONSUMER_MAXPROCESSINGTIME", "30")
	os.Setenv("CONSUMER_RETURN_ERROR", "false")
	os.Setenv("CONSUMER_OFFSETS_AUTO_COMMIT_ENABLE", "false")
	os.Setenv("CONSUMER_OFFSETS_AUTO_COMMIT_INTERVAL", "10")
	os.Setenv("CONSUMER_OFFSETS_INITIAL", "newest")
	os.Setenv("CONSUMER_OFFSETS_RETRY_MAX", "0")
	os.Setenv("CONSUMER_GROUP_SESSION_TIMEOUT", "50")
	os.Setenv("CONSUMER_GROUP_HEARTBEAT_INTERVAL", "30")
	os.Setenv("CONSUMER_GROUP_REBALANCE_STRATEGY", "rr")
	os.Setenv("CONSUMER_GROUP_REBALANCE_TIMEOUT", "10")
	os.Setenv("CONSUMER_GROUP_REBALANCE_TIMEOUT", "10")
	os.Setenv("CONSUMER_GROUP_REBALANCE_RETRY_MAX", "3")
	os.Setenv("CONSUMER_GROUP_REBALANCE_RETRY_BACKOFF", "3")
	os.Setenv("CONSUMER_CLIENTID", "test")
	os.Setenv("CONSUMER_CHANNELBUFFERSIZE", "3")
	os.Setenv("CONSUMER_VERSION", "2.7.0")
	os.Setenv("CONSUMER_BROKERS", "localhost:9092")
	os.Setenv("CONSUMER_GROUP", "testEpochConvertor")
	_, err := newConsumer()
	if err != nil {
		t.Error()
	}
}

func TestIncPrometheusCounter(t *testing.T) {
	err := incPrometuesCounter()
	if err != nil {
		t.Error()
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("CONSUMER_GROUP", "testEpochConvertor")
	if getEnv("CONSUMER_GROUP", "test") != "testEpochConvertor" {
		t.Error()
	}
}
