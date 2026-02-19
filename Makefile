build-1:
	@go build -C phase1-building-oidc-auth/ -o ../bin/runner
run-1: build-1
	@./bin/runner
