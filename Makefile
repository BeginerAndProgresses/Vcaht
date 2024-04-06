.PHONY: docker
docker:
	@rm v_chat || true
	@go mod tidy
	@set GOARCH=arm
	@go env -w GOARCH=arm
	@set GOOS=linux
	@go env -w GOOS=linux
	@go build -tags=k8s -o v_chat .
	@set GOARCH=amd64
	@go env -w GOARCH=amd64
	@set GOOS=windows
	@go env -w GOOS=windows
	@docker rmi -f ximubuqi/v_chat:v0.0.1
	@docker pull ubuntu:20.04
	@docker build -t ximubuqi/v_chat:v0.0.1 .