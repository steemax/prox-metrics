# Prox-metrics
    экспортер утилизации ресурсов из кластеров Proxmox (https://wikijs.wb.ru/ru/wbpayops/infra/prox-metrics)
    
 список метрик:
 ```
 # HELP node_cpu_alloc_counts Number of CPU core used (assigned) VMs on this node
# TYPE node_cpu_alloc_counts gauge
node_cpu_alloc_counts{node="test-proxmox-k8s-01"} 6
node_cpu_alloc_counts{node="test-proxmox-k8s-02"} 0
node_cpu_alloc_counts{node="test-proxmox-k8s-03"} 0
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

### Варианты запуска
 
1. Как сервис в Linux
2. Как сервис в Kubernetes

 ### Deploy "Как сервис в Linux"
 
 *предварительно настраиваем ./config/config.json, после*
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
*2 варианта установки ниже, HELM + K Apply*

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



## Проверка работы
Проверить можно вручную на линуксе просто запустив, 

```bash
cd ./deploy/linux
./prox-metrics
```

бинарь будет искать конфиг для коннекта к проксу, он должен лежать в том же каталоге где и бинарь по пути
./config/config.json
после запуска в stdout будет лог, по которому видно все ок или нет:

> 2023/12/11 11:42:12 Nodes Online:[test-proxmox-k8s-02 test-proxmox-k8s-03 test-proxmox-k8s-01]
2023/12/11 11:42:12 Avg CPU util for node test-proxmox-k8s-02: 0.049389
2023/12/11 11:42:12 CPU  count for node test-proxmox-k8s-02 is: 96
2023/12/11 11:42:12 CPU numbers assigned to VMs on node test-proxmox-k8s-02 is: 0
2023/12/11 11:42:12 MEM  count for node test-proxmox-k8s-02 is: 754 Gb
2023/12/11 11:42:12 MEM allocated size to VMs for node test-proxmox-k8s-02 is: 0 Gb
2023/12/11 11:42:12 Disk size allocated to VMs for node test-proxmox-k8s-02 is: 0 Gb
2023/12/11 11:42:12 Disk size for node test-proxmox-k8s-02 is: 1787 Gb
2023/12/11 11:42:12 Avg CPU util for node test-proxmox-k8s-03: 0.049731
2023/12/11 11:42:12 CPU  count for node test-proxmox-k8s-03 is: 96
2023/12/11 11:42:12 CPU numbers assigned to VMs on node test-proxmox-k8s-03 is: 0
2023/12/11 11:42:12 MEM  count for node test-proxmox-k8s-03 is: 754 Gb
2023/12/11 11:42:12 MEM allocated size to VMs for node test-proxmox-k8s-03 is: 0 Gb
2023/12/11 11:42:12 Disk size allocated to VMs for node test-proxmox-k8s-03 is: 0 Gb
2023/12/11 11:42:12 Disk size for node test-proxmox-k8s-03 is: 1787 Gb
2023/12/11 11:42:12 Avg CPU util for node test-proxmox-k8s-01: 0.051170
2023/12/11 11:42:12 CPU  count for node test-proxmox-k8s-01 is: 96
2023/12/11 11:42:12 CPU numbers assigned to VMs on node test-proxmox-k8s-01 is: 6
2023/12/11 11:42:12 MEM  count for node test-proxmox-k8s-01 is: 754 Gb
2023/12/11 11:42:12 MEM allocated size to VMs for node test-proxmox-k8s-01 is: 5 Gb
2023/12/11 11:42:13 Disk size allocated to VMs for node test-proxmox-k8s-01 is: 28 Gb
2023/12/11 11:42:13 Disk size for node test-proxmox-k8s-01 is: 1787 Gb