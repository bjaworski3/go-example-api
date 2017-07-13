VERSION ?= 0.0.3
APP_NAME ?= web-api
CLUSTER_NAME ?= web-api-cluster
PROJECT_ID ?= project-id
DOCKER_REG ?= us.gcr.io
GCP_ZONE ?= us-east1-b

DOCKER_IMAGE ?= $(DOCKER_REG)/$(PROJECT_ID)/$(APP_NAME):$(VERSION)
GOLANG_IMAGE ?= golang:1.8

all: help
	
help:
	 @echo " deploy          Deploy the go app as a Docker container to GKE"
	 @echo " undeploy        Undeploy the go app on GKE and clean"
	 @echo " test            Run the go tests"
	 @echo " clean           Remove all docker images for this app"
	 
deploy:
	docker build -t $(DOCKER_IMAGE) .
	gcloud docker -- push $(DOCKER_IMAGE)
	gcloud config set compute/zone $(GCP_ZONE)
	gcloud config set project $(PROJECT_ID)
	gcloud container clusters create $(CLUSTER_NAME)
	kubectl create -f deployment.yaml
	kubectl create -f web-api-service.yaml 
	@echo "Run 'kubectl get services web-api' until the external IP has been assigned"

undeploy: clean
	kubectl delete deployment $(APP_NAME)
	kubectl delete service $(APP_NAME)
	gcloud container clusters delete $(CLUSTER_NAME) -q

check test tests:
	go test
	
clean:
	-docker rmi $(DOCKER_IMAGE)
	-docker rmi $(GOLANG_IMAGE)