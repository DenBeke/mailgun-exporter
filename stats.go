package mailgunexporter

import "github.com/mailgun/mailgun-go/v4"

// Stats contain all the stats tracked in this exporter.
type Stats struct {
	Accepted          int
	Clicked           int
	Complained        int
	Delivered         int
	FailedPermanently int
	FailedTemporary   int
	Opened            int
	Stored            int
	Unsubscribed      int
}

// mailgunStatsToStats transforms the *mailgun.Stats object to our own *Stats object.
func mailgunStatsToStats(s *mailgun.Stats) *Stats {

	stats := Stats{
		Accepted:          s.Accepted.Total,
		Clicked:           s.Clicked.Total,
		Complained:        s.Complained.Total,
		Delivered:         s.Delivered.Total,
		FailedPermanently: s.Failed.Permanent.Total,
		FailedTemporary:   s.Failed.Temporary.Espblock,
		Opened:            s.Opened.Total,
		Stored:            s.Stored.Total,
		Unsubscribed:      s.Unsubscribed.Total,
	}

	return &stats

}
