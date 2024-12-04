# Defines what will be invoked whenever make is run without arguments
.DEFAULT_GOAL := up

# Modules
# each module corresponds to a .yml file in the project root directory
system_modules := \
 	auth \
	leaderboard \
	garage \
	racing \
	orchestrator

ifneq ($(service),)
	services = $(service)
endif

build:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) build $(services)

up:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) up $(services)

upd:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) up --detach $(services)

down:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) down $(services)

remove:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) rm -s $(services)

