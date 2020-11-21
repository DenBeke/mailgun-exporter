# Prometheus Exporter for Mailgun ✉️

[![Go Report Card](https://goreportcard.com/badge/github.com/DenBeke/mailgun-exporter)](https://goreportcard.com/report/github.com/DenBeke/mailgun-exporter)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/denbeke/mailgun-exporter?sort=date)](https://hub.docker.com/r/denbeke/mailgun-exporter)

WIP


## Run with Docker

You can use the [docker-compose.yml](./docker-compose.yml) file:

```yaml
version: "3"
services:
  mailgun-exporter:
    container_name: mailgunexporter
    image: denbeke/mailgun-exporter
    ports:
      - "9999:9999"
    environment:
      - MAILGUN_REGION=EU
      - MAILGUN_PRIVATE_API_KEY=<your-mailgun-key>
```


## Run in Development

    env MAILGUN_REGION=EU MAILGUN_PRIVATE_API_KEY=<your-mailgun-key> go run cmd/mailgunexporter/*.go



## Grafana Dashboard

A Grafana dashboard for these metrics is available here: [./grafana/dashboard.json](grafana/dashboard.json)

[![](grafana/dashboard.png)](grafana/dashboard.png)


## Author

[Mathias Beke](https://denbeke.be)