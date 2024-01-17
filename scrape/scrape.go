package scrape

import (
	"log"
	nodes "prox-metrics/api"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Create Prometheus Metric
	cpuUsageGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_cpu_avg_util",
			Help: "Average CPU usage (percent) of Proxmox nodes (SUM rrddata daily cpu from api / COUNT)",
		},
		[]string{"node"}, // Метки (label)
	)

	cpuCouuntGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_cpu_count",
			Help: "Number of CPU core for Proxmox node",
		},
		[]string{"node"}, // Метки (label)
	)

	cpuCouuntAssignGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_cpu_alloc_counts",
			Help: "Number of CPU core used (assigned) VMs on this node",
		},
		[]string{"node"}, // Метки (label)
	)

	memCouuntGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_mem_all",
			Help: "Memory Size (in Gb) for node",
		},
		[]string{"node"}, // Метки (label)
	)

	memCouuntAllocGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_mem_allocated",
			Help: "Memory Size (in Gb) assigned to VMs",
		},
		[]string{"node"}, // Метки (label)
	)

	diskSizeAllocGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_allocated",
			Help: "Disk Size (in Gb) assigned to VMs",
		},
		[]string{"node"}, // Метки (label)
	)

	diskSizeGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_total",
			Help: "Disk Size (in Gb) total",
		},
		[]string{"node"}, // Метки (label)
	)

	cpuCountPowerOffVmsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_cpu_alloc_pwroff_vms",
			Help: "Total CPUs assigned to VMs in power-off state",
		},
		[]string{"node"}, // Метки (label)
	)
)

func PromRegisterMetrics() {
	// Register Prometheus Metrics
	prometheus.MustRegister(cpuUsageGaugeVec, cpuCouuntGaugeVec, cpuCouuntAssignGaugeVec, memCouuntGaugeVec, memCouuntAllocGaugeVec, diskSizeAllocGaugeVec, diskSizeGaugeVec, cpuCountPowerOffVmsGaugeVec)
}

func ScrapeNodes(nodesList []string, client *nodes.Client) {
	for _, nodeName := range nodesList {
		avgCPU, err := nodes.GetCPUUsage(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node CPU daily utilization %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("Avg CPU util for node %s: %f\n", nodeName, avgCPU)

		cpuCount, err := nodes.GetCPUCount(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node CPU count %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("CPU  count for node %s is: %d\n", nodeName, cpuCount)

		cpuAlloc, err := nodes.GetCPUAlloc(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for number CPU assigned to VMs %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("CPU numbers assigned to VMs on node %s is: %d\n", nodeName, cpuAlloc)

		cpuAllocPwrOff, err := nodes.GetVmsOff(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for number CPU assigned to Power-off VMs %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("CPU numbers assigned to VMs in Power-off state, on node %s is: %d\n", nodeName, cpuAllocPwrOff)

		memCount, err := nodes.GetMemCount(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node MEM count %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("MEM  count for node %s is: %d Gb\n", nodeName, memCount)

		memCountAlloc, err := nodes.GetMemAlloc(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node MEM allocated %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("MEM allocated size to VMs for node %s is: %d Gb\n", nodeName, memCountAlloc)

		diskSizeAlloc, err := nodes.GetDiskAlloc(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node Disk allocated %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("Disk size allocated to VMs for node %s is: %d Gb\n", nodeName, diskSizeAlloc)

		diskSizeTot, err := nodes.GetDiskSize(client, nodeName)
		if err != nil {
			log.Printf("Error get response from API for node Disk size %s: %v\n", nodeName, err)
			continue
		}
		log.Printf("Disk size for node %s is: %d Gb\n", nodeName, diskSizeTot)

		// Set Prometheus Metric
		cpuUsageGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(avgCPU)
		cpuCouuntGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(cpuCount))
		cpuCouuntAssignGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(cpuAlloc))
		memCouuntGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(memCount))
		memCouuntAllocGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(memCountAlloc))
		diskSizeAllocGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(diskSizeAlloc))
		diskSizeGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(diskSizeTot))
		cpuCountPowerOffVmsGaugeVec.With(prometheus.Labels{"node": nodeName}).Set(float64(cpuAllocPwrOff))
	}
}
