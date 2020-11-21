package mailgunexporter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Serve(config *Config) {

	m, err := New(config.MailgunPrivateAPIKey, config.MailgunRegion)
	if err != nil {
		log.Fatalf("couldn't create mailgun exporter instance: %v", err)
	}

	e := echo.New()

	g := e.Group("/metrics")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err = m.CollectMetrics()
			if err != nil {
				log.Errorf("couldn't collect metrics: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "couldn't collect metrics from MailGun")
			}

			return next(c)
		}
	})

	g.GET("/", echo.WrapHandler(promhttp.Handler()))

	log.Fatal(e.Start(config.HTTPAddress))

}
