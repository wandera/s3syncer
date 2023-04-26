BIN = s3syncer

.PHONY: build
build: *.go **/*.go
	go build -o $(BIN)

.PHONY: clean
clean:
	$(RM) $(BIN)
