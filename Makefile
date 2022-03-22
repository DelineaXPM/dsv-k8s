# The name of the image which is not normally overridden
NAME:=dsv-injector

# The version is overridden to make tagged releases, e.g., v1.0.0
VERSION?=latest

# The Kubernetes Namespace that the webhook will be deployed in 📁
NAMESPACE?=dsv

# Your roles.json file; see the README.md
ROLES_JSON?=configs/roles.json

# 👇 Podman works too
DOCKER=docker

# Helm is required to install the webhook
HELM=helm

# The chart is in the 'charts' subfolder per Helm convention
HELM_CHART:=charts/$(NAME)

# Use the kubectl included with Minikube
KUBECTL=minikube kubectl --

# Get the location of the registry from Minikube; *it is assumed to be running* -- see the README.md ⚠️
REGISTRY=$(shell $(KUBECTL) get -n kube-system service registry -o jsonpath="{.spec.clusterIP}{':'}{.spec.ports[0].port}")

.PHONY: clean image install install-image release uninstall

all: install-image

# Build the dsv-injector service container image 📦
image:
	$(DOCKER) $(DOCKER_ARGS) build . $(DOCKER_BUILD_ARGS) -t $(NAME):$(VERSION)

# Publish the image to $(REGISTRY); Used by GitHub Actions
release: image
	$(DOCKER) $(DOCKER_ARGS) tag $(DOCKER_TAG_ARGS) $(NAME):$(VERSION) $(REGISTRY)/$(NAME):$(VERSION)
	$(DOCKER) $(DOCKER_ARGS) push $(DOCKER_PUSH_ARGS) $(REGISTRY)/$(NAME):$(VERSION)

# Install the Helm chart using a roles.json file 📄
install:
	$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) install --create-namespace \
	--set-file rolesJson=$(ROLES_JSON) $(HELM_INSTALL_ARGS) $(HELM_REPO_ARGS) \
	--wait $(NAME) $(HELM_CHART)
# Install the chart with the locally built image in place of the default ⚙️
install-image: HELM_REPO_ARGS = --set image.pullPolicy=Never,image.repository=$(NAME)
install-image: image install

# Uninstall the Helm Chart ❌
uninstall:
	-$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) uninstall $(NAME)

# Remove the Docker images 🗑️
clean:
	-$(DOCKER) $(DOCKER_ARGS) rmi -f $(NAME):$(VERSION)
