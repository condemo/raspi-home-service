binary-name=raspi_service

build: templ-gen
	@GOOS=windows GOARCH=amd64 go build -o ./bin/${binary-name}-win.exe ./cmd/api/main.go
	@GOOS=linux GOARCH=amd64 go build -o ./bin/${binary-name}-linux ./cmd/api/main.go
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/${binary-name}-darwin ./cmd/api/main.go

run: build
	@./bin/${binary-name}-linux

arm-build:
	@GOOS=linux GOARCH=arm64 go build -o ./bin/${binary-name}-arm64 ./cmd/api/main.go

arm-run: arm-build
	@./bin/${binary-name}-arm64

test:
	@go test cmd/api/main.go

clean:
	@rm -rf ./bin/*
	@go clean

css-build:
	@tailwindcss -i ./public/static/css/input.css -o ./public/static/css/style.css

css-watch:
	@tailwindcss -i ./public/static/css/input.css -o ./public/static/css/style.css --watch

templ-build:
	@templ generate

templ-watch:
	@templ generate --watch
