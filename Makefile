audit:
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...

tidy:
	go mod tidy -v
	go fmt ./...

