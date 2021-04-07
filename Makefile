default: clean download mod-tidy fmt vet-check test build

.PHONY: build clean clean-binary clean-frontend default download fmt mod-tidy test vet-check web-assets

build: trento
trento: web-assets
	go build

clean: clean-binary clean-frontend

clean-binary:
	go clean

clean-frontend:
	rm -rf web/frontend/assets
	rm -rf web/frontend/node_modules

download:
	go mod download
	go mod verify

fmt:
	go fmt ./...

mod-tidy:
	go mod tidy

test: download web-assets
	go test -v ./...

vet-check: download web-assets
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
