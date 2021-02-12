AZBI := epiphanyplatform/azbi:dev
AZKS := epiphanyplatform/azks:dev

#makes it easier to replace the value
M_NAME            ?= azks-module-tests
M_SHARED_EXTERNAL ?= /tmp/$(M_NAME)/shared

VMS_RSA_FILENAME ?= vms_rsa

define AZBI_CONFIG
{
	"kind": "azbi",
	"version": "v0.1.2",
	"params": {
		"name": "$(M_NAME)",
		"location": "northeurope",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "aks",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
endef

-include ./azure.mk

export

.PHONY: all apply destroy clean metadata setup

all: apply

apply: init-azbi plan-azbi apply-azbi \
       init-azks plan-azks apply-azks

destroy: destroy-plan-azks destroy-azks \
         destroy-plan-azbi destroy-azbi

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
		init
	@echo "$$AZBI_CONFIG" > $(M_SHARED_EXTERNAL)/azbi/azbi-config.json

plan-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		plan

apply-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		apply

output-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZBI) \
		output

destroy-plan-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		plan \
		--destroy

destroy-azbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		destroy

.PHONY: init-azks plan-azks apply-azks output-azks destroy-plan-azks destroy-azks

init-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		init

plan-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZKS) \
		plan

apply-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZKS) \
		apply

output-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AZKS) \
		output

destroy-plan-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZKS) \
		plan \
		--destroy

destroy-azks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZKS) \
		destroy

setup: $(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME)

$(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
