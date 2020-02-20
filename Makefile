NAME:=dsv-injector
VERSION?=latest

IMAGE_TAG=$(NAME):$(VERSION)

image:
	docker build . -t $(IMAGE_TAG) -f build/Dockerfile

### The recipes below build and deploy a test injector-svc 🥼🧪

# The CA certificate of the Kubernetes cluster 🔐
CA_CRT?=${HOME}/.minikube/ca.crt

# The TCP port on which the service should listen 🌐
SERVICE_PORT?=8543

# The Kubernetes cluster Namespace in which to deploy the Webhook and/or POD
NAMESPACE?=default

# The JSON file containing a mapping of DSV role names to tenant/credentials 🔑
ROLES_FILE?=configs/roles.json

# The IP address of the host running the dsv-injector service 🖥️
HOST_IP?=$(shell ip -o -4 -br addr | sed -n 2p | awk '{print $$3}' | sed -e 's|/.*$$||')

TEST_IMAGE_TAG?=$(NAME)-test:$(VERSION)

BUILD_DIR=target

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

ca_bundle=$(shell base64 -w0 $(CA_CRT))

deploy_webhook: $(BUILD_DIR)
	sed -e "s| namespace: .*$$| namespace: $(NAMESPACE)|" \
		-e "s| port: [0-9]*.*$$| port: $(SERVICE_PORT)|" \
		-e "s|caBundle:.*$$|caBundle: $(ca_bundle)|" \
		deployments/webhook.yml >| $(BUILD_DIR)/webhook.yml
	kubectl apply -f $(BUILD_DIR)/webhook.yml

# Create a test image that includes the certficate and roles.json 🔓😧

$(BUILD_DIR)/$(NAME).key $(BUILD_DIR)/$(NAME).pem: $(BUILD_DIR)
	sh scripts/get_cert.sh -n "$(NAME)" -N "$(NAMESPACE)" -d "$(BUILD_DIR)"

# The registry can be supplied; spin up a test registry when it is not 🏗️
registry:
ifndef REGISTRY
	kubectl apply -f test/registry.yml
	@echo "Waiting for the test registry to spin up..." && sleep 6
REGISTRY=$(shell kubectl get service -n kube-system kube-registry -o json |\
		jq -r ".spec.clusterIP,.spec.ports[0].port" | sed -e "N;s|\n|:|")
endif

$(BUILD_DIR)/Dockerfile: registry test/Dockerfile $(BUILD_DIR)
	sed -e "s|^FROM $(NAME):.*|FROM $(REGISTRY)/$(IMAGE_TAG)|" \
		test/Dockerfile >| $(BUILD_DIR)/Dockerfile

deploy_host: deploy_webhook $(BUILD_DIR)/$(NAME).key $(BUILD_DIR)/$(NAME).pem
	sed -e "s| namespace: .*$$| namespace: $(NAMESPACE)|" \
		-e "s|- port: [0-9]*.*$$|- port: $(SERVICE_PORT)|" \
		-e "s|- ip: *\"[0-9].*$$|- ip: \"$(HOST_IP)\"|" \
		deployments/host.yml >| $(BUILD_DIR)/host.yml
	kubectl apply -f $(BUILD_DIR)/host.yml

test_image: image $(BUILD_DIR)/$(NAME).key $(BUILD_DIR)/$(NAME).pem $(BUILD_DIR)/Dockerfile
	docker tag $(IMAGE_TAG) $(REGISTRY)/$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_TAG)
	docker build . -t $(TEST_IMAGE_TAG) -f $(BUILD_DIR)/Dockerfile \
		--build-arg cert_file="$(BUILD_DIR)/$(NAME).pem" \
		--build-arg key_file="$(BUILD_DIR)/$(NAME).key" \
		--build-arg roles_file="$(ROLES_FILE)"

deploy: deploy_webhook test_image
	sed -e "s| namespace: .*$$| namespace: $(NAMESPACE)|" \
		-e "s|- port: [0-9]*.*$$|- port: $(SERVICE_PORT)|" \
		-e "s|image:.*$$|image: $(TEST_IMAGE_TAG)|" \
		deployments/pod.yml >| $(BUILD_DIR)/pod.yml
	kubectl apply -f $(BUILD_DIR)/pod.yml

deploy_clean:
	docker rmi -f $(IMAGE_TAG)
	kubectl delete --ignore-not-found deployments $(NAME)
	kubectl delete --ignore-not-found service $(NAME)
	kubectl delete --ignore-not-found mutatingwebhookconfigurations.admissionregistration.k8s.io $(NAME)

clean: deploy_clean
	rm -rf $(BUILD_DIR) $(NAME)
