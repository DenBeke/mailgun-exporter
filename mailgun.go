package mailgunexporter

import (
	"context"
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

	acceptedGauge          *prometheus.GaugeVec
	deliveredGauge         *prometheus.GaugeVec
	failedTemporaryGauge   *prometheus.GaugeVec
	failedPermanentlyGauge *prometheus.GaugeVec
	openedGauge            *prometheus.GaugeVec
	clickedGauge           *prometheus.GaugeVec
	complainedGauge        *prometheus.GaugeVec
	unsubscribedGauge      *prometheus.GaugeVec
	storedGauge            *prometheus.GaugeVec
}

// New creates a new MailgunExoprter with the given private API key and region.
func New(privateAPIKey string, region string) (*MailgunExporter, error) {

	m := &MailgunExporter{
		privateAPIKey: privateAPIKey,
		region:        region,
	}

	labels := []string{"domain"}

	// TODO: add domain as label
	m.acceptedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_accepted_total",
		Help: "",
	}, labels)
	m.deliveredGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_delivered_total",
		Help: "",
	}, labels)
	m.failedTemporaryGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_failed_temporary_total",
		Help: "",
	}, labels)
	m.failedPermanentlyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_failed_permanently_total",
		Help: "",
	}, labels)
	m.openedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_opened_total",
		Help: "",
	}, labels)
	m.clickedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_clicked_total",
		Help: "",
	}, labels)
	m.complainedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_complained_total",
		Help: "",
	}, labels)
	m.unsubscribedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_unsubscribed_total",
		Help: "",
	}, labels)
	m.storedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mailgun_stored_total",
		Help: "",
	}, labels)

	return m, nil
}

// CollectMetrics will get all the stats for all the domains and put them in the correct prometheus collectors.
func (m *MailgunExporter) CollectMetrics() error {

	domains, err := m.ListDomains()
	if err != nil {
		log.Errorf("couldn't list domains: %v", err)
		return fmt.Errorf("couldn't list domains: %v", err)

	}
	log.Printf("domains: %+v", domains)

	for _, domain := range domains {

		stats, err := m.GetStats(domain)
		if err != nil {
			log.Errorf("couldn't get stats: %v", err)
			return fmt.Errorf("couldn't get stats: %v", err)
		}

		log.Printf("stats for %s: %+v", domain, stats)

		m.SetPrometheusFromStats(stats, domain)

	}

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

// SetPrometheusFromStats sets all the values from the stats object as values for the Prometheus gauges for the given domain.
func (m *MailgunExporter) SetPrometheusFromStats(stats *Stats, domain string) {
	// TODO: do this for all mailgun events

	labels := prometheus.Labels{
		"domain": domain,
	}

	m.acceptedGauge.With(labels).Set(float64(stats.Accepted))
	m.deliveredGauge.With(labels).Set(float64(stats.Delivered))
	m.failedTemporaryGauge.With(labels).Set(float64(stats.FailedTemporary))
	m.failedPermanentlyGauge.With(labels).Set(float64(stats.FailedPermanently))
	m.openedGauge.With(labels).Set(float64(stats.Opened))
	m.clickedGauge.With(labels).Set(float64(stats.Clicked))
	m.complainedGauge.With(labels).Set(float64(stats.Complained))
	m.unsubscribedGauge.With(labels).Set(float64(stats.Unsubscribed))
	m.storedGauge.With(labels).Set(float64(stats.Stored))
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
