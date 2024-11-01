# Auto Scale Sevice

## Introduction

[Video Presentation](https://www.youtube.com/watch?v=bv_r--E3l_8)

The system consists of an API and a worker, where the API is required to respond quickly to web requests and the worker is responsible for handling time-consuming tasks. They communicate with each other through message queues.

There are two methods here, one is a cronjob that simulates normal traffic requests. The other is a manual trigger that throws a large number of messages into the message queue to simulate high traffic.


```go
func (p *Service) triggerSendHugeAmountMessages(c context.Context) (err error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Minute)
	defer cancel()

	// Send a message every 10 millsecond, exit after 10 minutes
	every := 10 * time.Millisecond
	t := time.NewTicker(every)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
			p.Publish(
				context.Background(),
				def.Topics.DoTask,
				&model.DoTaskCommand{Name: "do task"},
			)
		}
	}
}

```

```go
func (p *Service) cronSendLittleMessages(c context.Context) (err error) {
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()

	// Send a message every 100 millsecond, exit after 1 minutes
	every := 100 * time.Millisecond
	t := time.NewTicker(every)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
			p.Publish(
				context.Background(),
				def.Topics.DoTask,
				&model.DoTaskCommand{Name: "do task"},
			)
		}
	}
}

```


`jobDoTask` consume those messages.

```
func (p *Service) jobDoTask(c context.Context, msg *pubsub.Message) {
	now := xtime.Now()
	var err error
	defer func() {
		msg.Ack()
		prom.Consumer.Timing(
			fmt.Sprintf("consumer:%s", def.Topics.DoTask),
			int64(time.Since(now)/time.Millisecond),
		)
		prom.Consumer.Incr(fmt.Sprintf("consumer:%s", def.Topics.DoTask))
	}()

	cmd := new(model.DoTaskCommand)
	if err = jsoniter.Unmarshal(msg.Data, cmd); err != nil {
		log.For(c).Errorf("jobSendMail error(%+v)", err)
		return
	}

	log.For(c).Info("Do some small task")
}

```



## Deploy Steps

### Install gcloud CLI 

* [Install the gcloud CLI](https://cloud.google.com/sdk/docs/install)
* `gcloud init`
* `gcloud auth application-default login`
* `gcloud components install kubectl`

### Provision Environment
```
terraform int 
terraform apply
```

### Config kubectl

```
gcloud container clusters get-credentials ascale-439911-gke --region asia-east2 --project ascale-439911
```

### Install Helm

[Install Helm3 ](https://helm.sh/docs/intro/install/)


### Build docker images and push to repository

* Go to source top directory

* Build & push
```
docker build -t gcr.io/ascale-439911/ascale-worker  ./ -f DockerfileWorker && docker push gcr.io/ascale-439911/ascale-worker && docker build -t gcr.io/ascale-439911/ascale  ./ -f Dockerfile && docker push gcr.io/ascale-439911/ascale
```
### Deply to GKE
* Crate Namespace `kubectl create ns uat`
* Add configmaps `kubectl -n=uat create configmap acesoconf --from-file=./config.toml`
* Add Gcloud Credentials `kubectl -n=uat create secret generic pubsub-key --from-file=key.json`
* Deploy App `helm -n=uat install ascale ./k8s`

### HPA Autoscale

 Install Keda 
```
cd 3rd/keda
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace keda --create-namespace
kubectl -n=uat create -f pubsub-scale.yaml
```
![HPA](./assets/hpa.png)

### Install Prometheus & Grafana
```
cd 3rd/prometheus
kubectl create ns monitoring
helm upgrade -i prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --values kube-prometheus-stack.values
kubectl -n=monitoring port-forward service/prometheus-grafana 8999:80
kubectl -n=monitoring get secret prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

#### QPS
![QPS](./assets/qps.png)

#### REPLICAS
![REPLICAS](./assets/replicas.png)

#### CPU & Memory
![CPU](./assets/cpu.png)

