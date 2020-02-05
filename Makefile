
NAME=dsv-injector
VERSION?=0.0.1

# The Kubernetes cluster Namespace in which to deploy the Webhook and/or POD
NAMESPACE?=default

# The JSON file containing a mapping of DSV role names to tenant/credential pairs
ROLES_FILE?=configs/roles.json

# The TCP port on which the service should listen
SERVICE_PORT?=8543

CA_CRT=${HOME}/.minikube/ca.crt

# The IP address of the host of the dsv-injector service
HOST_IP?=$(shell ip -o -4 -br addr show dev eth0 | awk '{print $$3}' | sed -e 's|/.*$$||')

BUILD_DIR?=.build
IMAGE_TAG?=$(NAME):v$(VERSION)

all: image deploy_host

build:
	mkdir -p $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)
	kubectl delete --ignore-not-found deployments $(NAME)
	kubectl delete --ignore-not-found service $(NAME)
	kubectl delete --ignore-not-found mutatingwebhookconfigurations.admissionregistration.k8s.io $(NAME)

cert=$(BUILD_DIR)/$(NAME).pem

$(cert): build
	sh scripts/get_cert.sh -n "$(NAME)" -N "$(NAMESPACE)" -d "$(BUILD_DIR)"

deploy_webhook: $(cert) image
	sed -e "s| port: [0-9]*.*$$| port: $(SERVICE_PORT)|" \
		-e "s|caBundle:.*$$|caBundle: $(shell base64 -w0 $(CA_CRT))|" \
		deployments/webhook.yml >| $(BUILD_DIR)/webhook.yml
	kubectl apply -f $(BUILD_DIR)/webhook.yml

deploy_host: deploy_webhook
	sed -e "s|- port: [0-9]*.*$$|- port: $(SERVICE_PORT)|" \
		-e "s|- ip: *\"[0-9].*$$|- ip: \"$(HOST_IP)\"|" \
		deployments/host.yml >| $(BUILD_DIR)/host.yml
	kubectl apply -f $(BUILD_DIR)/host.yml

deploy_pod: deploy_webhook
	sed -e "s|- port: [0-9]*.*$$|- port: $(SERVICE_PORT)|" \
		-e "s|image:.*$$|image: $(IMAGE_TAG)|" \
		deployments/pod.yml >| $(BUILD_DIR)/pod.yml
	kubectl apply -f $(BUILD_DIR)/pod.yml

image: $(cert)
	docker build -t $(IMAGE_TAG) . \
		--build-arg cert_file="$(BUILD_DIR)/$(NAME).pem" \
		--build-arg key_file="$(BUILD_DIR)/$(NAME).key" \
		--build-arg roles_file="$(ROLES_FILE)" \

