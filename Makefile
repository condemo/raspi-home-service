binary-name=default

build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/${binary-name}-win.exe ./cmd/api/main.go
	@GOOS=linux GOARCH=amd64 go build -o ./bin/${binary-name}-linux ./cmd/api/main.go
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/${binary-name}-darwin ./cmd/api/main.go

run: build
	@./bin/${binary-name}-linux


clean:
	@rm -rf ./bin/*
	@go clean

