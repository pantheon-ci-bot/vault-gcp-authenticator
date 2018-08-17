APP := vault-gcp-authenticator

deps: ## install deps
	@dep ensure

build: ## build binary for current architecture
	go build -o $(APP)

build-release: ## build linux binary for release
	@GOOS=linux GOARCH=amd64 go build -o $(APP)

push-release: build-release ## push a binary release to github releases
	@go get -u github.com/tcnksm/ghr
	@ghr -t $$GITHUB_TOKEN \
		-u $$CIRCLE_PROJECT_USERNAME \
		-r $$CIRCLE_PROJECT_REPONAME \
		"$$CIRCLE_BUILD_NUM" \
		./$(APP)

build-docker: ## build docker container
	@echo TODO

push-docker: ## push docker container
	@echo TODO

help: ## print list of tasks and descriptions
        @grep --no-filename -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?##"}; { printf "\033[36m%-30s\033[0m %s \n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: help
