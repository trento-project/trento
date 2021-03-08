default: clean download mod-tidy fmt vet-check test build

.PHONY: build clean clean-binary clean-frontend default download fmt mod-tidy test vet-check webapp-assets

build: console-for-sap-applications
console-for-sap-applications: webapp-assets
	go build

clean: clean-binary clean-frontend

clean-binary:
	go clean

clean-frontend:
	rm -rf webapp/frontend/assets
	rm -rf webapp/frontend/node_modules

download:
	go mod download
	go mod verify

fmt:
	go fmt ./...

mod-tidy:
	go mod tidy

test: download webapp-assets
	go test -v ./...

vet-check: download webapp-assets
	go vet ./...

webapp-deps: webapp/frontend/node_modules
webapp/frontend/node_modules:
	cd webapp/frontend; npm install

webapp-assets: webapp/frontend/assets

webapp/frontend/assets: webapp/frontend/assets/js webapp/frontend/assets/stylesheets webapp/frontend/assets/images

webapp/frontend/assets/js: webapp/frontend/node_modules
	mkdir -p webapp/frontend/assets/js/eos-ds
	cp webapp/frontend/javascripts/layout.js webapp/frontend/assets/js/layout.js
	cp webapp/frontend/node_modules/eos-ds/dist/js/index.js webapp/frontend/assets/js/eos-ds/index.js

webapp/frontend/assets/stylesheets: webapp/frontend/node_modules
	mkdir -p webapp/frontend/assets/stylesheets/eos-icons
	webapp/frontend/node_modules/.bin/sass \
		webapp/frontend/stylesheets/stylesheets.scss:webapp/frontend/assets/stylesheets/stylesheets.css
	cp webapp/frontend/node_modules/eos-ds/dist/vendors/eos-icons/css/eos-icons.css webapp/frontend/assets/stylesheets/eos-icons/eos-icons.css
	cp -R webapp/frontend/node_modules/eos-ds/dist/vendors/eos-icons/fonts webapp/frontend/assets/stylesheets/
	webapp/frontend/node_modules/.bin/sass \
		webapp/frontend/stylesheets/override.scss:webapp/frontend/assets/stylesheets/override.css

webapp/frontend/assets/images:
	mkdir -p webapp/frontend/assets/images
	cp -R webapp/frontend/images webapp/frontend/assets
