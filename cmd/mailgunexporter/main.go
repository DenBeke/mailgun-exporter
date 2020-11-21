package main

import (
	log "github.com/sirupsen/logrus"

	mailgunexporter "github.com/DenBeke/mailgun-exporter"
)

func main() {

	config := mailgunexporter.NewConfigFromEnv()
	err := config.Validate()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	mailgunexporter.Serve(config)

}
