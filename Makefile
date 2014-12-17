build:
	@go build -o goget

cover:
	@go test -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out
