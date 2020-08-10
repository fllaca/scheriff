_artifacts:
	mkdir -p _artifacts

.PHONY: build
build: _artifacts
	CGO_ENABLED=0 go build -o _artifacts/okay .

.PHONY: test
test: _artifacts
	go test -coverprofile=_artifacts/coverage.out -coverpkg github.com/fllaca/okay/... ./...
	go tool cover -func=_artifacts/coverage.out
	go tool cover -html=_artifacts/coverage.out -o _artifacts/coverage.html

