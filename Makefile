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

# Build all container images
build:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) build $(services)
	docker compose -f client/compose.yml build

# Run the whole system (local webserver for client and entire system)
up: up-system up-client

up-system:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) up $(services)

up-client:
	docker compose -f client/compose.yml up

# Run the whole system detached (local webserver for client and entire system)
upd: upd-system upd-client

upd-system:
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) up $(services) --detach

upd-client:
	docker compose -f client/compose.yml up --detach

# Down the whole system (local webserver for client and entire system)
down: down-system down-client
	docker compose $(foreach module,$(system_modules),-f system/$(module).yml) down $(services)
	docker compose -f client/compose.yml down
