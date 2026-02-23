build-1:
	@go build -C phase1-building-oidc-auth/ -o ../bin/runner
run-1: build-1
	@./bin/runner
build-2:
	@go build -C phase2-abstracting-auth-platform/ -o ../bin/runner
run-2: build-2
	@./bin/runner
