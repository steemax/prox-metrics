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
	"github.com/hashicorp/vault/api"
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
		Vault         bool          `json:"enable_vault"`
		VaultServer   string        `json:"vault_server"`
		VaultPath     string        `json:"vault_path"`
		VaultRoleID   string        `json:"vault_role_id"`
		VaultSecretID string        `json:"vault_secret_id"`
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

	// Vault client initialize
	var vaultClient *api.Client
	if config.Vault {
		log.Printf("Info: Initialize vault client...")
		vaultConfig := api.DefaultConfig()
		vaultConfig.Address = config.VaultServer
		vaultClient, err = api.NewClient(vaultConfig)
		if err != nil {
			panic(err)
		}

		// get token from AppRole
		loginData := map[string]interface{}{
			"role_id":   config.VaultRoleID,
			"secret_id": config.VaultSecretID,
		}

		resp, err := vaultClient.Logical().Write("auth/approle/login", loginData)
		if err != nil {
			log.Print("Error: approle/login: can not login to Vault use role_is and secret_id")
			panic(err)
		}

		log.Print("Info: approle/login: success get token from vault")
		vaultClient.SetToken(resp.Auth.ClientToken)

	}

	getVaultCredentials := func() (uid, secret string, err error) {
		secretValues, err := vaultClient.Logical().Read(config.VaultPath)
		if err != nil {
			return "", "", err
		}

		if secretValues == nil || secretValues.Data == nil {
			log.Print("Error: getVaultCredentials: secret not found Vault")
			return "", "", fmt.Errorf("secret not found in Vault")
		}

		data, ok := secretValues.Data["data"].(map[string]interface{})
		if !ok {
			log.Print("Error: getVaultCredentials: 'data' field not found in Vault secret")
			return "", "", fmt.Errorf("'data' field not found in Vault secret")
		}

		uid, ok = data["token_uid"].(string)
		if !ok {
			log.Print("Error: getVaultCredentials: token_uid not found in Vault secret")
			return "", "", fmt.Errorf("token_uid not found in Vault secret")
		}

		secret, ok = data["token_secret"].(string)
		if !ok {
			log.Print("Error: getVaultCredentials: token_secret not found in Vault secret")
			return "", "", fmt.Errorf("token_secret not found in Vault secret")
		}

		return uid, secret, nil
	}

	proxServer := config.ProxAddress
	proxPort := config.ProxPort
	proxScheme := config.ProxScheme
	scrapeTime := config.MetricTimeout * time.Second
	var tokenUid, tokenSecret string
	if config.Vault {
		tokenUid, tokenSecret, err = getVaultCredentials()
		if err != nil {
			panic(err)
		} else {
			log.Print("Info: success get secret for ProxAPI from vault")
			log.Print("Info: get Prox metrics use Vault data for API")
		}
	} else {
		log.Print("Info: Vault is disabled in config, get Prox metrics use local config data for API")
		tokenUid = config.ProxTokenUID
		tokenSecret = config.ProxTokenSec
	}

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
