DIST := ${CURDIR}/dist

.DEFAULT_GOAL := help
.PHONY: help
help:  ## show this help
	@grep -E '^[a-zA-Z_\/-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

########################################
# Build
########################################

.PHONY: build
build: build/upload ## build all commands

build/upload: build/%: FORCE
	go build -o ${DIST}/$* ${RACE} ${TAGS} ./cmd/$*

FORCE: # dummy target
