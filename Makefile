VERSION ?= $(shell ./hack/get_version_from_git.sh)
FLAVOR ?= "Community"
LDFLAGS = -X github.com/trento-project/trento/version.Version="$(VERSION)"
LDFLAGS := $(LDFLAGS) -X github.com/trento-project/trento/version.Flavor="$(FLAVOR)"
ARCHS ?= amd64 arm64 ppc64le s390x
DEBUG ?= 0

ifeq ($(DEBUG), 0)
	LDFLAGS += -s -w
	GO_BUILD = CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath
else
	GO_BUILD = CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"
endif

.PHONY: default
default: clean mod-tidy fmt vet-check web-check e2e-check test build

.PHONY: build
build: trento
trento: web/frontend/assets
	$(GO_BUILD)

.PHONY: cross-compiled $(ARCHS)
cross-compiled: $(ARCHS)
$(ARCHS): web-assets
	@mkdir -p build/$@
	GOOS=linux GOARCH=$@ $(GO_BUILD) -o build/$@/trento

.PHONY: clean
clean: clean-binary clean-frontend clean-e2e-tests

.PHONY: clean-binary
clean-binary:
	go clean
	rm -rf build

.PHONY: clean-frontend
clean-frontend: clean-web-assets clean-web-deps

.PHONY: clean-web-assets
clean-web-assets:
	rm -rf web/frontend/assets

.PHONY: clean-web-deps
clean-web-deps:
	rm -rf web/frontend/node_modules

.PHONY: clean-web-assets-js
clean-web-assets-js:
	rm -rf web/frontend/assets/js

.PHONY: clean-e2e-tests
clean-e2e-tests:
	rm -rf test/e2e/node_modules

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmt-check
fmt-check:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]

.PHONY: generate
generate:
ifeq (, $(shell command -v mockery 2> /dev/null))
	$(error "'mockery' command not found. You can install it locally with 'go install github.com/vektra/mockery/v2'.")
endif
ifeq (, $(shell command -v swag 2> /dev/null))
	$(error "'swag' command not found. You can install it locally with 'go install github.com/swaggo/swag/cmd/swag'.")
endif
	go generate ./...

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: test
test: web-assets
	GIN_MODE=test go test -v -p 1 ./...

.PHONY: full-check
full-check: generate vet-check test web-check e2e-check

.PHONY: test-coverage
test-coverage: build/coverage.out
build/coverage.out:
	@mkdir -p build
	GIN_MODE=test go test -cover -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out

.PHONY: vet-check
vet-check: web-assets
	go vet ./...

.PHONY: web-deps
web-deps: web/frontend/node_modules
web/frontend/node_modules:
	cd web/frontend; npm install

.PHONY: web-check
web-check: web-format-check web-lint

.PHONY: web-format
web-format:
	cd web/frontend; npx prettier --write .

.PHONY: web-format-check
web-format-check:
	cd web/frontend; npx prettier --check .

.PHONY: web-lint
web-lint:
	cd web/frontend; npx eslint .

.PHONY: e2e-deps
e2e-deps: test/e2e/node_modules
test/e2e/node_modules:
	cd test/e2e; npm install

.PHONY: e2e-check
e2e-check: e2e-deps e2e-format-check e2e-lint

.PHONY: e2e-format
e2e-format:
	cd test/e2e; npx prettier --write .

.PHONY: e2e-format-check
e2e-format-check:
	cd test/e2e; npx prettier --check .

.PHONY: e2e-lint
e2e-lint:
	cd test/e2e; npx eslint .

.PHONY: web-assets
web-assets: web/frontend/assets

web/frontend/assets: web/frontend/assets/js web/frontend/assets/stylesheets web/frontend/assets/images

web/frontend/assets/js: web/frontend/node_modules
	mkdir -p web/frontend/assets/js/eos-ds
	cp web/frontend/javascripts/*.js web/frontend/assets/js/
	cp web/frontend/node_modules/jquery/dist/jquery.min.js web/frontend/assets/js/
	cp web/frontend/node_modules/bootstrap/dist/js/bootstrap.bundle.min.js web/frontend/assets/js/
	cp web/frontend/node_modules/eos-ds/dist/js/index.js web/frontend/assets/js/eos-ds/index.js
	cp web/frontend/node_modules/eos-ds/dist/js/index.js web/frontend/assets/js/eos-ds/index.js
	cp web/frontend/node_modules/bootstrap-select/dist/js/bootstrap-select.min.js web/frontend/assets/js/
	cp web/frontend/node_modules/@yaireo/tagify/dist/tagify.min.js web/frontend/assets/js/
	cp web/frontend/node_modules/@yaireo/tagify/dist/tagify.polyfills.min.js web/frontend/assets/js/
	cd web/frontend; npx webpack

web/frontend/assets/stylesheets: web/frontend/node_modules
	mkdir -p web/frontend/assets/stylesheets/eos-icons
	web/frontend/node_modules/.bin/sass \
		web/frontend/stylesheets/stylesheets.scss:web/frontend/assets/stylesheets/stylesheets.css
	cp web/frontend/node_modules/eos-icons/dist/css/eos-icons.css web/frontend/assets/stylesheets/eos-icons/eos-icons.css
	cp web/frontend/node_modules/eos-icons/dist/css/eos-icons-outlined.css web/frontend/assets/stylesheets/eos-icons/eos-icons-outlined.css
	cp -R web/frontend/node_modules/eos-icons/dist/fonts web/frontend/assets/stylesheets/
	web/frontend/node_modules/.bin/sass \
		web/frontend/stylesheets/override.scss:web/frontend/assets/stylesheets/override.css
	cp web/frontend/node_modules/bootstrap/dist/css/bootstrap.min.css web/frontend/assets/stylesheets/
	cp web/frontend/node_modules/bootstrap/dist/css/bootstrap.min.css.map web/frontend/assets/stylesheets/
	cp web/frontend/node_modules/bootstrap-select/dist/css/bootstrap-select.min.css web/frontend/assets/stylesheets/
	cp web/frontend/node_modules/@yaireo/tagify/dist/tagify.css web/frontend/assets/stylesheets/

web/frontend/assets/images:
	mkdir -p web/frontend/assets/images
	cp -R web/frontend/images web/frontend/assets

.PHONY: helm-lint
helm-lint:
	docker run --rm -ti --name trento-chart-test -w /workdir -v $(shell pwd):/workdir quay.io/helmpack/chart-testing:v3.4.0 ct lint

.PHONY: install-wwhrd
install-wwhrd:
	which wwhrd || go install github.com/frapposelli/wwhrd@latest

.PHONY: check-licenses
check-licenses: install-wwhrd
	go mod vendor
	wwhrd check -q -f .wwhrd.yaml