    
metrics example (all metrics available on :9090/metrics):

 ```
 # HELP node_cpu_alloc_counts Number of CPU core used (assigned) VMs on this node
# TYPE node_cpu_alloc_counts gauge
node_cpu_alloc_counts{node="test-proxmox-k8s-01"} 6
node_cpu_alloc_counts{node="test-proxmox-k8s-02"} 0
node_cpu_alloc_counts{node="test-proxmox-k8s-03"} 0
# HELP node_cpu_alloc_pwroff_vms Total CPUs assigned to VMs in power-off state Name
# TYPE: node_cpu_alloc_pwroff_vms gauge
node_cpu_alloc_pwroff_vms{node="test-proxmox-k8s-01"} 0
node_cpu_alloc_pwroff_vms{node="test-proxmox-k8s-02"} 0
node_cpu_alloc_pwroff_vms{node="test-proxmox-k8s-03"} 0
# HELP node_cpu_avg_util Average CPU usage (percent) of Proxmox nodes (SUM rrddata daily cpu from api / COUNT)
# TYPE node_cpu_avg_util gauge
node_cpu_avg_util{node="test-proxmox-k8s-01"} 0.05116983912712025
node_cpu_avg_util{node="test-proxmox-k8s-02"} 0.04938919443876377
node_cpu_avg_util{node="test-proxmox-k8s-03"} 0.04973066893411988
# HELP node_cpu_count Number of CPU core for Proxmox node
# TYPE node_cpu_count gauge
node_cpu_count{node="test-proxmox-k8s-01"} 96
node_cpu_count{node="test-proxmox-k8s-02"} 96
node_cpu_count{node="test-proxmox-k8s-03"} 96
# HELP node_disk_allocated Disk Size (in Gb) assigned to VMs
# TYPE node_disk_allocated gauge
node_disk_allocated{node="test-proxmox-k8s-01"} 28
node_disk_allocated{node="test-proxmox-k8s-02"} 0
node_disk_allocated{node="test-proxmox-k8s-03"} 0
# HELP node_disk_total Disk Size (in Gb) total
# TYPE node_disk_total gauge
node_disk_total{node="test-proxmox-k8s-01"} 1787
node_disk_total{node="test-proxmox-k8s-02"} 1787
node_disk_total{node="test-proxmox-k8s-03"} 1787
# HELP node_mem_all Memory Size (in Gb) for node
# TYPE node_mem_all gauge
node_mem_all{node="test-proxmox-k8s-01"} 754
node_mem_all{node="test-proxmox-k8s-02"} 754
node_mem_all{node="test-proxmox-k8s-03"} 754
# HELP node_mem_allocated Memory Size (in Gb) assigned to VMs
# TYPE node_mem_allocated gauge
node_mem_allocated{node="test-proxmox-k8s-01"} 5
node_mem_allocated{node="test-proxmox-k8s-02"} 0
node_mem_allocated{node="test-proxmox-k8s-03"} 0
```

### Choose install option
 
1. as Linux service
2. as Kubernetes service

 ### Deploy "as Linux service"
 
 *first configure ./config/config.json, after*
```bash
cd ./deploy/linux
adduser --system --group --no-create-home prox-metrics
mkdir -p /usr/local/prox-metrics
cp -r ./config /usr/local/prox-metrics/
cp ./prox-metrics /usr/local/prox-metrics/prox-metrics
chown prox-metrics:prox-metrics /usr/local/prox-metrics/prox-metrics sudo chmod +x /usr/local/prox-metrics/prox-metrics
cp ./prox-metrics.service /etc/systemd/system/prox-metrics.service
systemctl daemon-reload
systemctl start prox-metrics.service
```
 
 ### Deploy "Как сервис в Kubernetes"
*HELM + K Apply*

Kubectl
```bash
cd ./deploy
kubectl apply -f ./cm.yaml
kubectl apply -f ./deployment.yaml
```
HELM
```bash
cd ./deploy/proxmox-metrics
helm install prox-metrics -f ./value.yaml ./ --namespace=monitoring
```

## Grafana Dashboard
Пример борды лежит в репе, в JSON формате который можно экспортнуть как темплейт (./deploy/grafana-dashboard.json)



## Check on linux
run manually:

```bash
cd ./deploy/linux
./prox-metrics
```
see stdout logs:

> 2023/12/11 11:42:12 Nodes Online:[test-proxmox-k8s-02 test-proxmox-k8s-03 test-proxmox-k8s-01]
2023/12/11 11:42:12 Avg CPU util for node test-proxmox-k8s-02: 0.049389
2023/12/11 11:42:12 CPU  count for node test-proxmox-k8s-02 is: 96
2023/12/11 11:42:12 CPU numbers assigned to VMs on node test-proxmox-k8s-02 is: 0
2023/12/11 11:42:12 MEM  count for node test-proxmox-k8s-02 is: 754 Gb
2023/12/11 11:42:12 MEM allocated size to VMs for node test-proxmox-k8s-02 is: 0 Gb
2023/12/11 11:42:12 Disk size allocated to VMs for node test-proxmox-k8s-02 is: 0 Gb
2023/12/11 11:42:12 Disk size for node test-proxmox-k8s-02 is: 1787 Gb

## Troubleshooting

If you see error after get online nodes:

>2024/01/10 14:29:25 Nodes Online:[px-01 px-03 px-02 px-04 px-05]
2024/01/10 14:29:25 Avg CPU util for node px-04: NaN
panic: runtime error: index out of range [0] with length 0

>goroutine 1 [running]:
prox-metrics/api.GetCPUCount(0x8d97f70, {0x8d26ae6, 0x5})

the access token does not have enough rights to access this information in the API