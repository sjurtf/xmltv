
.PHONY: build
build:
	docker build . -t xmltv:development

.PHONY: fmt
fmt:
	go run mvdan.cc/gofumpt -w ./

.PHONY: fmt
check:
	go run honnef.co/go/tools/cmd/staticcheck ./...
	go run golang.org/x/vuln/cmd/govulncheck -v ./...
