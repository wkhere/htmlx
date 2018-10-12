go:
	go fmt 		./...
	go test -cover	.
	go install	./...

benchmark:
	go test -bench=. .

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go benchmark cover
