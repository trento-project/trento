default: clean download mod-tidy fmt vet-check test build

build: webapp-assets console-for-sap
console-for-sap:
	go build

clean:
	go clean
	rm -rf webapp/frontend/assets/*.css
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

webapp-assets: webapp/frontend/assets/stylesheets.css
webapp/frontend/assets/stylesheets.css: webapp/frontend/node_modules
	webapp/frontend/node_modules/.bin/sass \
		webapp/frontend/stylesheets/stylesheets.scss:webapp/frontend/assets/stylesheets.css
