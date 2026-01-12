package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type engine struct {
	NCodes []string
	// A value in seconds.
	Cooldown int
	Webhook  string
	Client   tls_client.HttpClient
}

func createEngine() (*engine, error) {
	var engine engine

	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&engine)
	if err != nil {
		return nil, err
	}

	if engine.Cooldown <= 0 {
		engine.Cooldown = 1
	}

	dnsServers := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1"}
	dnsServer := dnsServers[rand.Intn(len(dnsServers))]

	jar := tls_client.NewCookieJar()
	dialer := net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5 * time.Second,
				}

				return d.DialContext(ctx, "udp", net.JoinHostPort(dnsServer, "53"))
			},
		},
	}

	clientOptions := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_117),
		tls_client.WithCookieJar(jar),
		tls_client.WithDialer(dialer),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewLogger(), clientOptions...)
	if err != nil {
		return nil, err
	}

	engine.Client = client

	return &engine, nil
}
