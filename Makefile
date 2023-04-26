BIN = s3syncer

.PHONY: build
build: *.go **/*.go
	go build -o $(BIN)

.PHONY: check
check:
	golangci-lint run

.PHONY: clean
clean:
	$(RM) $(BIN)
