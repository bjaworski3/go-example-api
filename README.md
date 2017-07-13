# Go Example Hello Web Service

This repo holds the code for a Go example Web API. This Go application can be deployed to Google Container Engine.

## Web Service Endpoints
Below is a table of the available endpoints for the Web API.
Any other route requested will return an HTTP 404 not found status.
If a route below is requested with a undefined method, it will return a 405 Invalid request method.

| Url | Request Method | Description | 
|---|---|---|
| /hello/:name | GET | Responds with `Hello, <name>!` | 
| /health | GET | Serves as a health check endpoint and responds with some system stats/metric | 
| /counts | GET | Returns a JSON response with the counts of how many times each name has been called, in the format below. `[ {"name": "alice", "count": 2}, {"name": "bob", "count": 1} ]` |
| /counts | DELETE | Resets the data so there are no counts (empty counts array) |


## Deploying the application to Google Container Engine

This section explains how to deploy the application to the project defined in the Makefile. It will create a cluster called `web-api-cluster` then create a deployment and service called `web-api`.

* Need to login first:
`gcloud auth login`

* Run the deployment target `make deploy`. This will:
  * Build a local docker image
  * Upload the docker image to Google container registry
  * Create the `web-api-cluster` cluster
  * Create the `web-api` deployment
  * Create a `web-api` service for forwarding traffic to the container
  

* After the commands are finished, find the application external IP with: `kubectl get services web-api`

### Undeploy application
Delete the cluster and application with `make undeploy`. This will also delete the local docker images. The docker image will still be available on the Google container registry.

## Development

This section describes how to run the application locally 

### Required Packages

unix (gopsutil requirement):

`go get -v golang.org/x/sys/unix`

gopsutil:

`go get -v github.com/shirou/gopsutil`

### Running the application locally
* Build application:
`go install`

* Run the application
`go-example-api`

* Access the api locally at: <http://127.0.0.1:8080>

### Running the tests
* Run tests:
`go test`
