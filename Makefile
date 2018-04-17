NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
APP=pairtermd
REVISION=$(shell git rev-parse --short HEAD)
BASE_VERSION=$(shell cat VERSION)
VERSION=$(BASE_VERSION)-$(REVISION)

all: build

build:
	@echo "$(OK_COLOR)==> Building revision $(VERSION)...$(NO_COLOR)"
	@script/build $(APP) $(VERSION) dev

run: build
	@echo "$(OK_COLOR)==> Running revision $(VERSION)...$(NO_COLOR)"
	@script/run

run-client:
	@echo "$(OK_COLOR)==> Running pairtermjs revision $(VERSION)...$(NO_COLOR)"
	@script/run-client

build-client:
	@echo "$(OK_COLOR)==> Building pairtermjs for production revision $(VERSION)...$(NO_COLOR)"
	@script/build-client

prod: build-client rice
	@echo "$(OK_COLOR)==> Bundling assets in rice $(VERSION)...$(NO_COLOR)"
	@script/build $(APP) $(VERSION) prod

rice:
	@echo "$(OK_COLOR)==> Bundling assets in rice $(VERSION)...$(NO_COLOR)"
	@script/rice

yarn:
	@echo "$(OK_COLOR)==> Running yarn install revision $(VERSION)...$(NO_COLOR)"
	@script/yarn

format:
	go fmt ./...

test:
	@echo "$(OK_COLOR)==> Testing...$(NO_COLOR)"
	@script/test $(TEST)

release:
	@script/release $(VERSION)

.PHONY: all build test release run-client build-client
