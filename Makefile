NAME:=dsv-injector
HELM_CHART:=charts/$(NAME)
VERSION?=latest

# The Kubernetes Namespace that the webhook will be deployed in 📁
NAMESPACE?=default

# Your roles.json file; see the README.md)
ROLES_JSON?=configs/roles.json

# The CA certificate (chain); we are assuming Minikube; Minishift is similar. 💡
CA_BUNDLE?=${HOME}/.minikube/ca.crt

# 👇 Podman works too
DOCKER=docker

# Helm is required to install the webhook
HELM=helm

.PHONY: image uninstall docker-rmi remove-cert clean

all: install

# Build the dsv-injector service container image 📦
image:
	$(DOCKER) build . -t $(NAME):$(VERSION) $(DOCKER_BUILD_ARGS)

# Unless it already exists, get a certificate from the Kubernetes cluster CA 🔐
$(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem:
	sh scripts/get_cert.sh -n "$(NAME)" -d "$(HELM_CHART)" -N "$(NAMESPACE)"
	-rm -f $(HELM_CHART)/$(NAME).csr

# Install will use the cert and key below, no matter how they got there. 😉😇
install: $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem image
	$(HELM) install $(HELM_INSTALL_ARGS) \
	--set-file caBundle=$(CA_BUNDLE),rolesJson=$(ROLES_JSON) \
	--set image.repository=$(NAME),image.tag=$(VERSION) \
	$(NAME) $(HELM_CHART)

# Uninstall the Helm Chart and remove the Docker images
uninstall:
	-$(HELM) uninstall $(NAME)

# Remove the Docker images
docker-rmi:
	-$(DOCKER) rmi -f $(NAME):$(VERSION)

# Remove the X.509 certificate and RSA private key
remove-cert:
	-rm -f $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem

clean: docker-rmi remove-cert uninstall
