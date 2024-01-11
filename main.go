package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	nodes "prox-metrics/api"
	scrape "prox-metrics/scrape"

	"github.com/go-acme/lego/v4/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	type Config struct {
		ProxAddress   string        `json:"server"`
		ProxPort      string        `json:"server_port"`
		ProxScheme    string        `json:"protocol"`
		ProxTokenUID  string        `json:"token_uid"`
		ProxTokenSec  string        `json:"token_secret"`
		MetricTimeout time.Duration `json:"scrape_timeout"`
	}

	// Read the configuration file
	configApp, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}

	// Unmarshal the configuration data
	var config Config
	if err := json.Unmarshal(configApp, &config); err != nil {
		panic(err)
	}

	proxServer := config.ProxAddress
	proxPort := config.ProxPort
	proxScheme := config.ProxScheme
	tokenUid := config.ProxTokenUID
	tokenSecret := config.ProxTokenSec
	scrapeTime := config.MetricTimeout * time.Second

	scrape.PromRegisterMetrics()         // Initialize Prom metrics
	ticker := time.NewTicker(scrapeTime) // time duration from config for scrape metrics
	defer ticker.Stop()

	// for scrape after start app immediately
	runNow := make(chan struct{}, 1)
	runNow <- struct{}{}

	// Create URL
	url := fmt.Sprintf("%s://%s:%s/api2/json", proxScheme, proxServer, proxPort)

	// Create Token
	token := fmt.Sprintf("PVEAPIToken=%s=%s", tokenUid, tokenSecret)

	// Create Client
	client := nodes.NewClient(url, token)

	nodesList, err := client.GetOnlineNodes()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	log.Print("Nodes Online:", nodesList)

	// HTTP server for prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9090", nil) // metrics port

	for {
		select {
		case <-ticker.C:
			scrape.ScrapeNodes(nodesList, client) // scrape by timeout
		case <-runNow:
			scrape.ScrapeNodes(nodesList, client) //scrape by start app
		}
	}
}
