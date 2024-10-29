# Install Prometheuse

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack
```



## QPS
```
sum by(name) (rate(go_consumer_count{name="consumer:uat-do-task"}[10s]))
```
## Replicas 
```
kube_deployment_status_replicas{namespace="uat", deployment="ascale-worker"}
```
