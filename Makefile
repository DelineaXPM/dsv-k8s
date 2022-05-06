INJECTOR_NAME:=dsv-injector
SYNCER_NAME:=dsv-syncer

# The name of the image which is not normally overridden
IMAGE_NAME:=dsv-k8s

# The version is overridden to make tagged releases, e.g., v1.0.0
VERSION?=latest

# The Kubernetes Namespace that the webhook will be deployed in 📁
NAMESPACE?=dsv

# The default registry to push to and install from
REGISTRY:=quay.io/delinea

# Your credentials.json file; see the README.md
CREDS_JSON?=configs/credentials.json

# 👇 Podman works too
DOCKER=docker

# Helm is required to install the webhook
HELM=helm

# The chart is in the 'charts' subfolder per Helm convention
HELM_CHARTS:=./charts

# Supply the credentials.json file
HELM_INSTALL_ARGS = --set-file credentialsJson=$(CREDS_JSON)

# Use the kubectl included with Minikube
KUBECTL=kubectl

.PHONY: clean-image clean-injector clean-syncer clean image install-injector
	install-syncer install-image install-host install release uninstall-injector
	uninstall-syncer uninstall

all: binaries image

ifeq ($(OS),Windows_NT)
INJECTOR_BIN:=$(INJECTOR_NAME).exe
SYNCER_BIN:=$(SYNCER_NAME).exe
else
INJECTOR_BIN:=$(INJECTOR_NAME)
SYNCER_BIN:=$(SYNCER_NAME)
endif

$(INJECTOR_BIN): cmd/injector/main.go
	go build -o $@ $(GO_BUILD_FLAGS) ./cmd/injector

$(SYNCER_BIN): cmd/syncer/main.go
	go build -o $@ $(GO_BUILD_FLAGS) ./cmd/syncer

# build the platform-dependent binaries (for debugging) 🔨
binaries: $(INJECTOR_BIN) $(SYNCER_BIN)

# test the binaries 🧪
test:
	go test $(GO_TEST_FLAGS) ./pkg/...

# Build the dsv-injector service container image 📦
image:
	$(DOCKER) $(DOCKER_ARGS) build . $(DOCKER_BUILD_ARGS) -t $(IMAGE_NAME):$(VERSION)

# Tage the image with the version number 🎯
tag: image
	$(DOCKER) $(DOCKER_ARGS) tag $(DOCKER_TAG_ARGS) $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

# Publish the image to $(REGISTRY); Used by GitHub Actions 📢
release: tag
	$(DOCKER) $(DOCKER_ARGS) push $(DOCKER_PUSH_ARGS) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

# Install the webhook into Kubernetes 📛
install-injector:
	$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) install --create-namespace $(HELM_INSTALL_ARGS) $(HELM_REPO_ARGS) $(INJECTOR_NAME) $(HELM_CHARTS)/$(INJECTOR_NAME)

# Install the syncer into Kubernetes ⏲️
install-syncer:
	$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) install --create-namespace $(HELM_INSTALL_ARGS) $(HELM_REPO_ARGS) $(SYNCER_NAME) $(HELM_CHARTS)/$(SYNCER_NAME)

# Install via $(REGISTRY) presumably targeting a remote cluster 🌎
install: HELM_REPO_ARGS = --set image.repository=$(REGISTRY)/$(IMAGE_NAME)
install: install-syncer install-injector

# Install the webhook so that it points to the host; there is no deployment; this is for debugging 🥼🥽
install-host: CA_BUNDLE_KUBE_CONFIG_INDEX = 0
install-host: CA_BUNDLE_JSON_PATH = {.clusters[$(CA_BUNDLE_KUBE_CONFIG_INDEX)].cluster.certificate-authority-data}
install-host: CA_BUNDLE=$(shell $(KUBECTL) config view --raw -o jsonpath='$(CA_BUNDLE_JSON_PATH)' | tr -d '"')
install-host: HELM_REPO_ARGS = --set service.type=ExternalName,externalName=$(EXTERNAL_NAME),caBundle=$(CA_BUNDLE)
install-host: $(INJECTOR_BIN) install-injector

# Install with the locally built image; use this for Docker Desktop and Minikube ⚙️
install-image: HELM_REPO_ARGS = --set image.pullPolicy=Never,image.repository=$(IMAGE_NAME)
install-image: image install-syncer install-injector

# Push the locally built image to a remote registry then install it from there 📡
install-cluster: HELM_REPO_ARGS = --set imagePullPolicy=always,image.repository=$(REGISTRY)/$(IMAGE_NAME)
install-cluster: release image install-syncer install-injector

uninstall-injector:
	-$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) uninstall $(INJECTOR_NAME)

uninstall-syncer:
	-$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) uninstall $(SYNCER_NAME)

# Uninstall both Helm Charts ❌
uninstall: uninstall-syncer uninstall-injector

# Remove the Docker images
clean-image:
	-$(DOCKER) $(DOCKER_ARGS) rmi -f $(REGISTRY)/$(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):$(VERSION)

clean-injector:
	-rm -f $(INJECTOR_BIN)

clean-syncer:
	-rm -f $(SYNCER_BIN)

# Delete the image and binaries 🧹🗑️
clean: clean-image clean-injector clean-syncer
