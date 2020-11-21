# Prometheus Exporter for Mailgun ✉️

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


## Author

[Mathias Beke](https://denbeke.be)