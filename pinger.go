package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Pinger
type Pinger struct {
	Interval   time.Duration
	DestURL    url.URL
	httpClient *http.Client
	cancel     func()
}

func NewPinger(every time.Duration, destURL string) (*Pinger, error) {
	u, err := url.Parse(destURL)
	if err != nil {
		return nil, err
	}

	return &Pinger{Interval: every, DestURL: *u, httpClient: http.DefaultClient}, nil
}

func (p Pinger) Start(ctx context.Context) {
	log.Printf("Starting pinger for url: %s every %s", p.DestURL.String(), p.Interval.String())
	pingCycle := time.NewTicker(p.Interval)

	for {
		select {
		case <-pingCycle.C:
			if err := p.ping(ctx); err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("stopping pinger")
			return
		}
	}
}

func (p Pinger) ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.DestURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("got response with status", resp.Status)

	return nil
}

func (p *Pinger) Stop() {
	p.cancel()
}
