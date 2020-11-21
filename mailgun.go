package mailgunexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	log "github.com/sirupsen/logrus"
)

// MailgunExporter contains everything we need (including the Prometheus collectors)
type MailgunExporter struct {
	privateAPIKey string
	region        string

	acceptedGauge          prometheus.Gauge
	deliveredGauge         prometheus.Gauge
	failedTemporaryGauge   prometheus.Gauge
	failedPermanentlyGauge prometheus.Gauge
}

// New creates a new MailgunExoprter with the given private API key and region.
func New(privateAPIKey string, region string) (*MailgunExporter, error) {

	m := &MailgunExporter{
		privateAPIKey: privateAPIKey,
		region:        region,
	}

	// TODO: add domain as label
	m.acceptedGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mailgun_accepted_total",
		Help: "",
	})
	m.deliveredGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mailgun_delivered_total",
		Help: "",
	})
	m.failedTemporaryGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mailgun_failed_temporary_total",
		Help: "",
	})
	m.failedPermanentlyGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mailgun_failed_permanently_total",
		Help: "",
	})

	return m, nil
}

// CollectMetrics will get all the stats for all the domains and put them in the correct prometheus collectors.
func (m *MailgunExporter) CollectMetrics() error {

	_, err := m.ListDomains()
	if err != nil {
		log.Errorf("couldn't list domains: %v", err)
		return fmt.Errorf("couldn't list domains: %v", err)

	}
	//fmt.Printf("%+v", domains)

	stats, err := m.GetStats("denbeke.be")
	if err != nil {
		log.Errorf("couldn't get stats: %v", err)
		return fmt.Errorf("couldn't get stats: %v", err)
	}

	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		log.Errorf("couldn't marshal stats: %v", err)
		return fmt.Errorf("couldn't marshal stats: %v", err)
	}

	fmt.Println(string(jsonBytes))

	m.SetPrometheusFromStats(stats)

	return nil
}

// createMailgunAPIClient creates a Mailgun API client for the given domain and the current region.
func (m *MailgunExporter) createMailgunAPIClient(domain string) *mailgun.MailgunImpl {
	mg := mailgun.NewMailgun(domain, m.privateAPIKey)

	if strings.ToUpper(m.region) == "EU" {
		mg.SetAPIBase(mailgun.APIBaseEU)
	}

	return mg
}

// SetPrometheusFromStats sets all the values from the stats object as values for the Prometheus gauges.
func (m *MailgunExporter) SetPrometheusFromStats(stats *Stats) {
	// TODO: do this for all domains and add domain as label
	// TODO: do this for all mailgun events
	m.acceptedGauge.Set(float64(stats.Accepted))
	m.deliveredGauge.Set(float64(stats.Delivered))
	m.failedTemporaryGauge.Set(float64(stats.FailedTemporary))
	m.failedPermanentlyGauge.Set(float64(stats.FailedPermanently))
}

// GetStats returns the Mailgun stats for a given domain.
func (m *MailgunExporter) GetStats(domain string) (*Stats, error) {

	mg := m.createMailgunAPIClient(domain)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	stats, err := mg.GetStats(ctx, []string{
		"accepted",
		"clicked",
		"complained",
		"delivered",
		"failed",
		"opened",
		"stored",
		"unsubscribed",
	}, &mailgun.GetStatOptions{
		Duration:   "1d",
		Resolution: "day",
	})
	if err != nil {
		return nil, err
	}
	if len(stats) != 1 {
		return nil, fmt.Errorf("expected exactly one range of stats from API. got %d", len(stats))
	}

	return mailgunStatsToStats(&stats[0]), nil
}

// ListDomains returns all the domains in the current region.
func (m *MailgunExporter) ListDomains() ([]string, error) {

	mg := m.createMailgunAPIClient("")

	it := mg.ListDomains(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Domain
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}

	domains := []string{}

	for _, domain := range result {
		domains = append(domains, domain.Name)
	}

	return domains, nil
}
