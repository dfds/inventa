module github.com/dfds/inventa/operator

go 1.15

require (
	github.com/aws/aws-sdk-go v1.37.17
	github.com/coreos/go-oidc/v3 v3.1.0
	github.com/getoutreach/goql v1.4.0
	github.com/go-logr/logr v0.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/net v0.0.0-20210908191846-a5e095526f91 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
)
