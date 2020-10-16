certbot-driver: $(shell find . -type f -name *.go)
	CGO_ENABLED=0 go build -o "$@" ./cmd/certbot-driver
	@if ! ldd "$@" 2> /dev/null; then echo "OK: not a dynamic executable!"; fi

.PHONY: clean
clean:
	rm -Rfv certbot-bot

.PHONY: cl
cl:
	find . -type f -name *.go | xargs wc -l

.PHONY: test
test:
	go test ./...
