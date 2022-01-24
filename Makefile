# make のみで実行された場合 help を実行
.DEFAULT_GOAL := help


.PHONY: run
run: ## デーモン化せずローカルで実行
	go run ./cmd --daemonize=false


.PHONY: docker-build
docker-build: ## Docker イメージのビルド
	bash docker-build.sh


# Self-Documented Makefile
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
