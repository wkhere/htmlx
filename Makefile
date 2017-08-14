go:
	go fmt 		./...
	go vet		./...
	go build	./...
	go test -cover 	./...
	go install	./...

benchmark:
	go test -bench=. ./...

cover:
	go test -coverprofile=cov.find.out	./htmlx/find
	@echo now paste to run: go tool cover -html cov...
	@echo go tool cover -html cov. | xsel

.PHONY: go benchmark cover
