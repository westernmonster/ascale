# Auto Scale Sevice


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
gcloud container clusters get-credentials ascale-gke --zone asia-east2-a --project ascale-439911
```

### Install Helm

[Install Helm3 ](https://helm.sh/docs/intro/install/)


### Build docker images and push to repository

* Go to source top directory

```
docker build -t gcr.io/ascale-439911/ascale-worker  ./ -f DockerfileWorker && docker push gcr.io/ascale-439911/ascale-worker && docker build -t gcr.io/ascale-439911/ascale  ./ -f Dockerfile && docker push gcr.io/ascale-439911/ascale
```

* helm -n=uat install ascale ./k8s
