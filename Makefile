go:
	go vet 		./...
	go test -cover	.
	go build	./cmd/...

install: go
	go install -ldflags=-s ./cmd/...

bench:
	go test -bench=$(sel) -count $(cnt) -benchmem
sel=.
cnt=5

cover:
	go test -coverprofile=cov
	go tool cover -html=cov -o cov.html && browse cov.html

.PHONY: go install benchmark cover
