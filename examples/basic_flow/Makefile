AZBI := epiphanyplatform/azbi:0.0.1
AZKS := epiphanyplatform/azks:0.0.1

#makes it easier to replace the value
M_NAME            ?= mop-module-tests
M_SHARED_EXTERNAL ?= /tmp/$(M_NAME)/shared

VMS_RSA_FILENAME ?= vms_rsa

-include ./azure.mk

export

# replace the values if You want to run AKS in already existing subnet
EXISTING_RG_NAME ?= unset
EXISTING_SUBNET_ID ?= unset

.PHONY: all apply destroy clean metadata setup

all: apply

apply: init-azbi plan-azbi apply-azbi \
       init-azks plan-azks apply-azks

apply-only-azks: init-azks plan-azks apply-azks

destroy: destroy-plan-azks destroy-azks \
         destroy-plan-azbi destroy-azbi

destroy-only-azks: destroy-plan-azks destroy-azks

clean:
	@rm -rf $(M_SHARED_EXTERNAL)

metadata: guard-VERSION guard-IMAGE guard-USER
	@docker run --rm \
		-t $(USER)/$(IMAGE):$(VERSION) \
		metadata

.PHONY: init-azbi plan-azbi apply-azbi output-azbi destroy-plan-azbi destroy-azbi

init-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		init \
		M_VMS_COUNT=2 \
		M_PUBLIC_IPS=false \
		M_NAME=$(M_NAME)

plan-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		plan \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

apply-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		apply \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

output-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		output \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-plan-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		plan-destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

.PHONY: init-azks plan-azks apply-azks output-azks destroy-plan-azks destroy-azks

init-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		init \
		M_VMS_COUNT=2 \
		M_PUBLIC_IPS=false \
		M_NAME=$(M_NAME) \
        M_VMS_RSA=$(VMS_RSA_FILENAME) \
        M_EXISTING_RG_NAME=$(EXISTING_RG_NAME) \
        M_EXISTING_SUBNET_ID=$(EXISTING_SUBNET_ID)

plan-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		plan \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

apply-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		apply \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID \
		kubeconfig

output-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		output \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-plan-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		plan-destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

setup: $(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME)

$(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
