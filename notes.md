
`make docker-build docker-push IMG="example.com/memcached-operator:v0.0.1"`
error building at STEP "RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go": error while running runtime: exit status 1

`make bundle IMG="example.com/memcached-operator:v0.0.1"`
panic: runtime error: invalid memory address or nil pointer dereference 
make: *** [Makefile:100: manifests] Error 2

`operator-sdk create api --group=http --version=v1alpha1 --kind=HTTPServer`
FATA[0010] failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v4": exit status 2 