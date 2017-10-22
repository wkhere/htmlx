go:
	go fmt 		./...
	go vet		./...
	go build	./...
	go test -cover	.
	go install	./...

dep:
	go get golang.org/x/net/html

benchmark:
	go test -bench=. .

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go dep benchmark cover
