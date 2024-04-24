# Setup

### Prerequites

Operator SDK CLI: `brew install operator-sdk`

[go version 1.21](https://rpmfind.net/linux/rpm2html/search.php?query=golang&submit=Search+...&system=&arch=)

[golang-src 1.21](https://rpmfind.net/linux/rpm2html/search.php?query=golang-src&submit=Search+...&system=&arch=)

[golang-bin 1.21](https://rpmfind.net/linux/rpm2html/search.php?query=golang-bin&submit=Search+...&system=&arch=)

[go-filesystem](https://rpmfind.net/linux/rpm2html/search.php?query=go-filesystem&submit=Search+...&system=&arch=)

[git](https://git-scm.com/)

[docker](https://docs.docker.com/get-docker/) 

[kubectl & kind or minikube](https://kubernetes.io/docs/tasks/tools/)

`kind create cluster`

### What was done

Created project in http-operator:
`operator-sdk init --domain=example.com --repo=github.com/example/http-server-operator`

Create API (must be on golang 1.21):
`operator-sdk create api --group=http --version=v1alpha1 --kind=HTTPServer`

Updated the spec of api/v1alpha1/httpserver_types.go with the necessary fields for configuring the HTTP server.

Updated the controllers/httpserver_controller.go file to define the reconcile logic for the HTTP servers. This logic will manage the lifecycle of the HTTP servers based on the CRD specifications.

### Build and Deploy

#### HTTP Server

`cd src/basic_servers`

`docker build -t <your-docker-repo>/<image-name>:<tag> .`

`kind load docker-image <image-name>:<tag>`

Find and replace "your-http-server-image" with your image.

#### Operator

Run the operator:

`make manifests`

`make install run`

In a new terminal tab/window, create a sample HTTP custom resource:

`cd src/http-operator`

`kubectl apply -f config/crd/bases/http.example.com_httpservers.yaml`


Verify the Deployment:

```bash
kubectl get pods -n <namespace>
kubectl get deployments -n <namespace>
kubectl get service -n <namespace>
```

Test the Operator:

Create instances of your CRD and observe if the Operator creates, updates, or deletes resources as expected based on the CRD specifications.

Monitor and Troubleshoot:

Monitor the logs of the Operator to check for any errors or unexpected behavior.

`kubectl logs <operator-pod-name> -n <namespace>`
