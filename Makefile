NAME:=dsv-injector
HELM_CHART:=charts/$(NAME)
VERSION?=latest

# The Kubernetes Namespace that the webhook will be deployed in 📁
NAMESPACE?=default

# Your roles.json file; see the README.md)
ROLES_JSON?=configs/roles.json

# 👇 Podman works too
DOCKER=docker

# Helm is required to install the webhook
HELM=helm

.PHONY: cert clean clean-docker clean-cert image install install-image uninstall

all: install

# Build the dsv-injector service container image 📦
image:
	$(DOCKER) $(DOCKER_ARGS) build . -t $(NAME):$(VERSION) $(DOCKER_BUILD_ARGS)

# Create a self-signed SSL certificate 🔐
$(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem:
	sh scripts/get_cert.sh -n "$(NAME)" -d "$(HELM_CHART)" -N "$(NAMESPACE)"
cert: $(HELM_CHART)/$(NAME).pem

# Install will use the cert and key below, no matter how they got there. 😉😇
install: $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem
	$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) install --create-namespace \
	--set-file caBundle=$(HELM_CHART)/$(NAME).pem,rolesJson=$(ROLES_JSON) \
	$(HELM_INSTALL_ARGS) $(HELM_REPO_ARGS) $(NAME) $(HELM_CHART)
# Install image uses the locally built image in place of the default ⚙️
install-image: HELM_REPO_ARGS = --set image.pullPolicy=Never,image.repository=$(NAME)
install-image: image install

# Uninstall the Helm Chart
uninstall:
	-$(HELM) $(HELM_ARGS) --namespace $(NAMESPACE) uninstall $(NAME)

# Remove the X.509 certificate and RSA private key
clean-cert:
	-rm -f $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem

# Remove the Docker images
clean-docker:
	-$(DOCKER) $(DOCKER_ARGS) rmi -f $(NAME):$(VERSION)

clean: clean-cert clean-docker
