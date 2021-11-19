package main

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type PingerType int

const (
	website PingerType = iota
	api
	bot
)

// Pinger
type Pinger struct {
	Interval   time.Duration
	DestURL    url.URL
	httpClient *retryablehttp.Client
	cancel     func()
	Type       PingerType
}

func NewPinger(every time.Duration, destURL string) (*Pinger, error) {
	u, err := url.Parse(destURL)
	if err != nil {
		return nil, err
	}

	retryClient := retryablehttp.NewClient()
	return &Pinger{Interval: every, DestURL: *u, httpClient: retryClient, Type: website}, nil
}

func (p Pinger) Start(ctx context.Context) {
	log.Printf("Starting pinger for url: %s every %s", p.DestURL.String(), p.Interval.String())
	pingCycle := time.NewTicker(p.Interval)

	for {
		select {
		case <-pingCycle.C:
			switch p.Type {
			case website:
				if err := p.ping(ctx); err != nil {
					log.Println(err)
				}
			default:
				log.Println("unsupported PingerType")
			}
		case <-ctx.Done():
			log.Println("stopping pinger")
			return
		}
	}
}

func (p Pinger) ping(ctx context.Context) error {
	req, err := retryablehttp.NewRequest("GET", p.DestURL.String(), nil)
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
