go:
	go vet 		./...
	go test -cover	.
	go install	./...

bench:
	go test -bench=$(sel) -count $(cnt) -benchmem
sel=.
cnt=5

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go benchmark cover
