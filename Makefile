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

#### RUN ####

build: update_proto
	docker compose --profile run $(foreach module,$(system_modules),-f system/$(module).yml) build $(services)

up:
	docker compose --profile run $(foreach module,$(system_modules),-f system/$(module).yml) up $(services)

upd:
	docker compose --profile run $(foreach module,$(system_modules),-f system/$(module).yml) up --detach $(services)

stop:
	docker compose --profile run $(foreach module,$(system_modules),-f system/$(module).yml) stop $(services)

down:
	docker compose --profile run $(foreach module,$(system_modules),-f system/$(module).yml) down --volumes $(services)

#### TEST ####

build_test: down_test update_proto
	docker compose --profile test $(foreach module,$(system_modules),-f system/$(module).yml) build $(foreach service,$(services), test_$(service))

test: build_test
	docker compose --profile test $(foreach module,$(system_modules),-f system/$(module).yml) up $(foreach service,$(services), test_$(service))

stop_test:
	docker compose --profile test $(foreach module,$(system_modules),-f system/$(module).yml) stop

down_test:
	docker compose --profile test $(foreach module,$(system_modules),-f system/$(module).yml) down --volumes

#### SETUP and CLEANUP ####

clean: stop stop_test down down_test remove_compiler

update_proto: mkdir_proto compile_protobuf
	$(foreach module,$(system_modules),cp -r system/proto/. system/$(module)/proto ;)

compile_protobuf: pull_compiler
	docker run -v $(shell pwd)/system/proto:/defs namely/protoc-all -f service.proto -l go -o . --go-source-relative  

pull_compiler:
	docker pull namely/protoc-all

remove_compiler:
	docker rmi -f namely/protoc-all

mkdir_proto:
	mkdir -p $(foreach module,$(system_modules), system/$(module)/proto)

