package metrics

import (
	"os"
	"testing"
	"time"

	"github.com/lifeforms/httpcheck/httpcheck"
)

func TestExporter(t *testing.T) {
	os.Setenv("LISTEN_PORT", "8080")
	os.Setenv("METRICS_PATH", "/metrics")
	produce := ProducedMessageCounter.WithLabelValues("0", "output")
	consume := ConvertorMessageCounter.WithLabelValues("100", "input")

	produce.Inc()
	consume.Inc()
	go Exporter()
	time.Sleep(2 * time.Second)
	httpcheck.Verbose = true
	httpcheck.RequestTimeout = 5
	httpcheck.ServerTimeout = 10

	manifest := httpcheck.Manifest{
		httpcheck.Server{
			Name: "exporter",
			Scenarios: []httpcheck.Scenario{
				httpcheck.Scenario{
					Name:  "produced metrics",
					Tests: httpcheck.Tests{httpcheck.Test{Url: "http://localhost:8080/metrics", Content: "produced_message_total"}}},
				httpcheck.Scenario{
					Name:  "consumed metrics",
					Tests: httpcheck.Tests{httpcheck.Test{Url: "http://localhost:8080/metrics", Content: "convert_message_total"}}},
			},
		},
	}

	if manifest.Test() != nil {
		t.Error()
	}

}
