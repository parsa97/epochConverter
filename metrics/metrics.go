package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const exporterTemplate = `<html>
<head><title>Producer Exporter</title></head>
<body>
<p><a href= %s >Metrics</a></p>
</body>
</html>
`

var ProducedMessageCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "produced_message_total",
		Help: "How many message produced in which partition.",
	},
	[]string{"partition", "topic"},
)
var ConvertorMessageCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "convert_message_total",
		Help: "How many message converted.",
	},
	[]string{"consumerID", "topic"},
)

func Exporter(errCH chan error) {
	log.Info("Server Start Exporter")
	listernPort, metricsPath := exporterConfig()
	prometheus.MustRegister(ProducedMessageCounter)
	prometheus.MustRegister(ConvertorMessageCounter)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(exporterTemplate, metricsPath)))
	})
	http.Handle(metricsPath, promhttp.Handler())
	err := http.ListenAndServe(":"+listernPort, nil)
	if err != nil {
		errCH <- err
	}
}
