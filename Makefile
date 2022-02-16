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
	$(DOCKER) build . -t $(NAME):$(VERSION) $(DOCKER_BUILD_ARGS)

cert: $(HELM_CHART)/$(NAME).pem

install-image: image
	make install HELM_INSTALL_ARGS="--set image.repository=$(NAME)"

# Create a self-signed SSL certificate 🔐
$(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem:
	sh scripts/get_cert.sh -n "$(NAME)" -d "$(HELM_CHART)" -N "$(NAMESPACE)"

# Install will use the cert and key below, no matter how they got there. 😉😇
install: $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem
	$(HELM) install $(HELM_INSTALL_ARGS) \
	--set-file caBundle=$(HELM_CHART)/$(NAME).pem,rolesJson=$(ROLES_JSON) \
	$(NAME) $(HELM_CHART)

# Uninstall the Helm Chart and remove the Docker images
uninstall:
	-$(HELM) uninstall $(NAME)

# Remove the Docker images
clean-docker:
	-$(DOCKER) rmi -f $(NAME):$(VERSION)

# Remove the X.509 certificate and RSA private key
clean-cert:
	-rm -f $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem

clean: clean-docker clean-cert
