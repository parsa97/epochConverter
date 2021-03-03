package metrics

import "os"

func exporterConfig() (listenPort string, metricsPath string) {
	listenPort = getEnv("LISTEN_PORT", "8081")
	metricsPath = getEnv("METRICS_PATH", "/metrics")
	return listenPort, metricsPath
}

// GetEnv - Allows us to supply a fallback option if nothing specified
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
