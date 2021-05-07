sel=. # selection for bench

go:
	go fmt 		./...
	go test -cover	.
	go install	./...

bench:
	go test -bench=$(sel) -benchmem .

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go benchmark cover
