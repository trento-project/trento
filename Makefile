VERSION ?= $(shell hack/get_version_from_git.sh)
LDFLAGS = -X github.com/trento-project/trento/version.Version="$(VERSION)"
ARCHS ?= amd64 arm64 ppc64le s390x
DEBUG ?= 0

ifeq ($(DEBUG), 0)
	LDFLAGS += -s -w
	GO_BUILD = CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath
else
	GO_BUILD = CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"
endif

default: clean mod-tidy fmt vet-check test build

.PHONY: build clean clean-binary clean-frontend cross-compiled default fmt fmt-check generate swag mod-tidy test vet-check web-assets

build: trento
trento: web-assets
	$(GO_BUILD)

cross-compiled: $(ARCHS)
$(ARCHS): web-assets
	@mkdir -p build
	GOOS=linux GOARCH=$@ $(GO_BUILD) -o build/trento-$@

clean: clean-binary clean-frontend

clean-binary:
	go clean
	rm -rf build

clean-frontend:
	rm -rf web/frontend/assets
	rm -rf web/frontend/node_modules

fmt:
	go fmt ./...

fmt-check:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]

generate:
ifeq (, $(shell command -v mockery 2> /dev/null))
	$(error "'mockery' command not found. You can install it locally with 'go install github.com/vektra/mockery/v2'.")
endif
	go generate ./...

swag:
ifeq (, $(shell command -v swag 2> /dev/null))
	$(error "'swag' command not found. You can install it locally with 'go install github.com/swaggo/swag/cmd/swag'.")
endif
	swag init

mod-tidy:
	go mod tidy

test: generate web-assets
	go test -v ./...

test-coverage:
	@mkdir -p build
	go test -cover -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out

vet-check: generate web-assets
	go vet ./...

web-deps: web/frontend/node_modules
web/frontend/node_modules:
	cd web/frontend; npm install

web-assets: web/frontend/assets

web/frontend/assets: web/frontend/assets/js web/frontend/assets/stylesheets web/frontend/assets/images

web/frontend/assets/js: web/frontend/node_modules
	mkdir -p web/frontend/assets/js/eos-ds
	cp web/frontend/javascripts/*.js web/frontend/assets/js/
	cp web/frontend/node_modules/eos-ds/dist/js/index.js web/frontend/assets/js/eos-ds/index.js

web/frontend/assets/stylesheets: web/frontend/node_modules
	mkdir -p web/frontend/assets/stylesheets/eos-icons
	web/frontend/node_modules/.bin/sass \
		web/frontend/stylesheets/stylesheets.scss:web/frontend/assets/stylesheets/stylesheets.css
	cp web/frontend/node_modules/eos-ds/dist/vendors/eos-icons/css/eos-icons.css web/frontend/assets/stylesheets/eos-icons/eos-icons.css
	cp -R web/frontend/node_modules/eos-ds/dist/vendors/eos-icons/fonts web/frontend/assets/stylesheets/
	web/frontend/node_modules/.bin/sass \
		web/frontend/stylesheets/override.scss:web/frontend/assets/stylesheets/override.css

web/frontend/assets/images:
	mkdir -p web/frontend/assets/images
	cp -R web/frontend/images web/frontend/assets
