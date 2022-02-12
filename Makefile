NAME:=dsv-injector
VERSION?=latest

IMAGE_TAG=$(NAME):$(VERSION)

HELM_CHART:=charts/$(NAME)

DOCKER=docker
HELM=helm

# The Kubernetes Namespace in which to deploy 📁
NAMESPACE?=default

CA_BUNDLE?=${HOME}/.minikube/ca.crt

ROLES_JSON?=configs/roles.json

.PHONY: image

all: install

# Build the dsv-injector service image 📦
image:
	$(DOCKER) build . -t $(IMAGE_TAG) $(DOCKER_BUILD_ARGS)

# Get a certificate from the Kubernetes cluster CA
$(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem:
	sh scripts/get_cert.sh -n "$(NAME)" -d "$(HELM_CHART)" -N "$(NAMESPACE)"
	-rm -f $(HELM_CHART)/$(NAME).csr

install: $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem image
	$(HELM) install $(HELM_INSTALL_ARGS) \
	--set-file caBundle=$(CA_BUNDLE) \
	--set-file rolesJson=$(ROLES_JSON) \
	$(NAME) $(HELM_CHART)

clean:
	$(HELM) uninstall $(NAME)
	$(DOCKER) rmi -f $(IMAGE_TAG)
	-rm -f $(HELM_CHART)/$(NAME).key $(HELM_CHART)/$(NAME).pem
