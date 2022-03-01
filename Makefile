NAME:=dsv-injector
HELM_CHART:=charts/$(NAME)
VERSION?=latest

# The Kubernetes Namespace that the webhook will be deployed in 📁
NAMESPACE?=dsv

# Your roles.json file; see the README.md
ROLES_JSON?=configs/roles.json

# 👇 Podman works too
DOCKER=docker

# Helm is required to install the webhook
HELM=helm

.PHONY: clean image install install-image uninstall

all: install

# Build the dsv-injector service container image 📦
image:
	$(DOCKER) $(DOCKER_ARGS) build . $(DOCKER_BUILD_ARGS) -t $(NAME):$(VERSION)

# Publish the image to $(REGISTRY)
release: image
	$(DOCKER) $(DOCKER_ARGS) tag $(DOCKER_TAG_ARGS) $(NAME):$(VERSION) $(REGISTRY)/$(NAME):$(VERSION)
	$(DOCKER) $(DOCKER_ARGS) push $(DOCKER_PUSH_ARGS) $(REGISTRY)/$(NAME):$(VERSION)

# Install the Helm chart using a roles.json file 📄
install:
	$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) install --create-namespace \
	--set-file rolesJson=$(ROLES_JSON) $(HELM_INSTALL_ARGS) $(HELM_REPO_ARGS) \
	$(NAME) $(HELM_CHART)
# Install the chart with the locally built image in place of the default ⚙️
install-image: HELM_REPO_ARGS = --set image.pullPolicy=Never,image.repository=$(NAME)
install-image: image install

# Uninstall the Helm Chart ❌
uninstall:
	-$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) uninstall $(NAME)

# Remove the Docker images 🗑️
clean:
	-$(DOCKER) $(DOCKER_ARGS) rmi -f $(NAME):$(VERSION)
