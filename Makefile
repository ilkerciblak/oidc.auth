build-1:
	@go build -C phase1-building-oidc-auth/ -o ../bin/runner
run-1: build-1
	@./bin/runner
build-2:
	@go build -C phase2-abstracting-auth-platform/ -o ../bin/runner
run-2: build-2
	@./bin/runner
build-3:
	@go build -C phase3-implementing-packages/ -o ../bin/runner
run-3: build-3
	@./bin/runner
