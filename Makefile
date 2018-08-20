APP := vault-gcp-authenticator
IMAGE := quay.io/getpantheon/$(APP)

ifdef CIRCLE_BUILD_NUM
	TAG := $$CIRCLE_BUILD_NUM
else
	TAG := dev
endif

deps: ## install deps
	@dep ensure

build: ## build binary for current architecture
  @CGO_ENABLED=0 go build -o $(APP)

build-release: ## build linux binary for release
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP)

push-release: build-release ## push a binary release to github releases
	@go get -u github.com/tcnksm/ghr
	@ghr -t $$GITHUB_TOKEN \
		-u $$CIRCLE_PROJECT_USERNAME \
		-r $$CIRCLE_PROJECT_REPONAME \
		"$$CIRCLE_BUILD_NUM" \
		./$(APP)

build-docker: build-release ## build docker container
	@docker build -t $(IMAGE):$(TAG) .

push-docker: ## push docker container
	@docker push $(IMAGE):$(TAG)

help: ## print list of tasks and descriptions
	@grep --no-filename -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?##"}; { printf "\033[36m%-30s\033[0m %s \n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: help all deps build build-release push-release build-docker push-docker
