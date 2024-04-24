
`make docker-build docker-push IMG="example.com/memcached-operator:v0.0.1"`
error building at STEP "RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go": error while running runtime: exit status 1

`make bundle IMG="example.com/memcached-operator:v0.0.1"`
panic: runtime error: invalid memory address or nil pointer dereference 
make: *** [Makefile:100: manifests] Error 2

`operator-sdk create api --group=http --version=v1alpha1 --kind=HTTPServer`
FATA[0010] failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v4": exit status 2 

https://github.com/operator-framework/operator-sdk/issues/6681

must be on go version 1.21

`make manifests`
Error: not all generators ran successfully
run `controller-gen rbac:roleName=manager-role crd webhook paths=./... output:crd:artifacts:config=config/crd/bases -w` to see all available markers, or `controller-gen rbac:roleName=manager-role crd webhook paths=./... output:crd:artifacts:config=config/crd/bases -h` for usage
make: *** [Makefile:100: manifests] Error 1