package nodes

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Client API Proxmox
type Client struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// NewClient
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// GetOnlineNodes
func (c *Client) GetOnlineNodes() ([]string, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/nodes", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Data []struct {
			Node   string `json:"node"`
			Status string `json:"status"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var nodes []string
	for _, n := range data.Data {
		if n.Status == "online" {
			nodes = append(nodes, n.Node)
		}
	}

	return nodes, nil
}

func GetCPUUsage(client *Client, nodeName string) (float64, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/rrddata?timeframe=day", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			CPU float64 `json:"cpu"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	// Calculate avg CPU Util per node
	var sum float64
	for _, data := range result.Data {
		sum += data.CPU
	}
	avg := sum / float64(len(result.Data)) * 100
	//avg = math.Round(avg*1000) / 1000 // окруление, 3х знаков после запятой

	return avg, nil
}

func GetCPUCount(client *Client, nodeName string) (int, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/rrddata?timeframe=day", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			MaxCPU int `json:"maxcpu"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Data[0].MaxCPU, nil
}

func GetMemCount(client *Client, nodeName string) (int64, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/rrddata?timeframe=day", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check for a successful response status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// JSON Response decode
	var result struct {
		Data []struct {
			MaxMEM int64 `json:"memtotal"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	// Check if there is at least one data entry
	if len(result.Data) == 0 {
		return 0, errors.New("No data entries found in the API response")
	}

	// Check if MaxMEM is a valid number
	if result.Data[0].MaxMEM <= 0 {
		return 0, errors.New("Invalid MaxMEM value in the API response")
	}

	mem := result.Data[0].MaxMEM / 1024 / 1024 / 1024
	return mem, nil
}

func GetCPUAlloc(client *Client, nodeName string) (int, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/qemu", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			AllCPU int `json:"cpus"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	// Calculate total CPU assign to VMs
	var sum int
	for _, data := range result.Data {
		sum += data.AllCPU
	}

	return sum, nil
}

func GetMemAlloc(client *Client, nodeName string) (int64, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/qemu", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			AllMEM int64 `json:"maxmem"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	// Calculate total CPU assign to VMs
	var sum int64
	for _, data := range result.Data {
		sum += data.AllMEM
	}
	memAlloc := sum / 1024 / 1024 / 1024
	return memAlloc, nil
}

func GetDiskSize(client *Client, nodeName string) (int64, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/storage", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			Storage string `json:"storage"`
			Disk    int64  `json:"total"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	// Find the total disk size for storage "local"
	var disk int64
	for _, storage := range result.Data {
		if storage.Storage == "local" {
			disk = storage.Disk / 1024 / 1024 / 1024
			break
		}
	}

	if disk == 0 {
		return 0, errors.New("No storage data found for 'local' in the response")
	}

	return disk, nil
}

func GetDiskAlloc(client *Client, nodeName string) (int64, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/qemu", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			AllDisk int64 `json:"maxdisk"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	// Calculate total CPU assign to VMs
	var sum int64
	for _, data := range result.Data {
		sum += data.AllDisk
	}
	diskAlloc := sum / 1024 / 1024 / 1024
	return diskAlloc, nil
}

// Get VM in power-off state
func GetVmsOff(client *Client, nodeName string) (int, error) {
	// Create URL
	url := fmt.Sprintf("%s/nodes/%s/qemu", client.BaseURL, nodeName)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set Header with token
	req.Header.Set("Authorization", client.Token)

	// Get data from API using the embedded http.Client
	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// JSON Response decode
	var result struct {
		Data []struct {
			Status string `json:"status"`
			CPUs   int    `json:"cpus"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	// Calculate total CPUs assigned to VMs in power-off state
	var totalCPUs int = 0
	for _, vm := range result.Data {
		if vm.Status == "stopped" {
			totalCPUs += vm.CPUs
		}
	}

	return totalCPUs, nil
}
