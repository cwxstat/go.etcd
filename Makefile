
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)



##@ Tools

.PHONY: tools
tools: ## All software: KinD, Kubectl, etcdctl and etcd.
	@./.scripts/install-etcd.sh  
	@./.scripts/install-kind.sh  
	@./.scripts/install-kubectl.sh


##@ All targets.

.PHONY: complete-kind
complete-kind: kind cert-manager load  ## All - does everything


##@ Steps

.PHONY: kind 
kind: ## Build Kubernetes cluster named  'test'
	kind delete cluster --name=test
	kind create cluster --name=test
	@./.scripts/wait-for-cluster.sh


.PHONY: load  
load: ## load pod and service definitions
	@kubectl create configmap etcd-config --from-file=./certs/
	@k apply -f pod.yaml
	@sleep 14
	@kubectl port-forward pods/etcd 2379:2379 -n default &

.PHONY: cert-manager 
cert-manager: ## cert-manager v1.7.2
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.7.2/cert-manager.yaml


##@ etcd 

.PHONY: etcd-for-testing 
etcd-for-testing: ## etcd for testing... assums running pod 
	@go run main.go  config
	@etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem  --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem user add root --interactive=false --password=A08auslkdjMMf
