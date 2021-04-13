default: clean mod-tidy fmt vet-check test build

.PHONY: build clean clean-binary clean-frontend default fmt generate mod-tidy test vet-check web-assets

build: trento
trento: web-assets
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w'

clean: clean-binary clean-frontend

clean-binary:
	go clean

clean-frontend:
	rm -rf web/frontend/assets
	rm -rf web/frontend/node_modules

fmt:
	go fmt ./...

generate:
ifeq (, $(shell command -v mockery 2> /dev/null))
	$(error "'mockery' command not found. You can install it locally with 'go get github.com/vektra/mockery/v2'.")
endif
	go generate ./...

mod-tidy:
	go mod tidy

test: generate web-assets
	go test -v ./...

vet-check: web-assets
	go vet ./...

web-deps: web/frontend/node_modules
web/frontend/node_modules:
	cd web/frontend; npm install

web-assets: web/frontend/assets

web/frontend/assets: web/frontend/assets/js web/frontend/assets/stylesheets web/frontend/assets/images

web/frontend/assets/js: web/frontend/node_modules
	mkdir -p web/frontend/assets/js/eos-ds
	cp web/frontend/javascripts/layout.js web/frontend/assets/js/layout.js
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
