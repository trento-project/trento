default: clean download mod-tidy fmt vet-check test build

build: webapp-assets console-for-sap
console-for-sap:
	go build

clean:
	go clean
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

webapp-assets: webapp-assets-stylesheets webapp-assets-javascrips

webapp-assets-stylesheets: stylesheets stylesheets-eos-icons stylesheets-override
	mkdir -p webapp/frontend/assets/stylesheets

stylesheets: webapp/frontend/node_modules
	webapp/frontend/node_modules/.bin/sass \
		webapp/frontend/stylesheets/stylesheets.scss:webapp/frontend/assets/stylesheets/stylesheets.css

stylesheets-eos-icons: webapp/frontend/node_modules
	mkdir -p webapp/frontend/assets/stylesheets/eos-icons/
	cp webapp/frontend/node_modules/eos-ds/dist/vendors/eos-icons/css/eos-icons.css webapp/frontend/assets/stylesheets/eos-icons/eos-icons.css
	cp -R webapp/frontend/node_modules/eos-ds/dist/vendors/eos-icons/fonts webapp/frontend/assets/stylesheets/

stylesheets-override: webapp/frontend/node_modules
	webapp/frontend/node_modules/.bin/sass \
		webapp/frontend/stylesheets/override.scss:webapp/frontend/assets/stylesheets/override.css

webapp-assets-javascrips: javascripts-eos-ds
	mkdir -p webapp/frontend/assets/js

javascripts-eos-ds: webapp/frontend/node_modules
	mkdir -p webapp/frontend/assets/js/eos-ds
	cp webapp/frontend/node_modules/eos-ds/dist/js/index.js webapp/frontend/assets/js/eos-ds/index.js
