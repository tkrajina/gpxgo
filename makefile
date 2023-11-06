.PHONY: test
test: compile
	go test ./gpx

.PHONY: compile
compile:
	go build -gcflags "-e" -o /dev/null ./...

.PHONY: gofmt
gofmt:
	gofmt -w ./gpx

.PHONY: goimports
goimports:
	goimports -w ./gpx

.PHONY: install
install:
	go install ./gpx

.PHONY: prepare
prepare:
	go get

.PHONY: clean
clean:
	echo "TODO"

.PHONY: ctags
ctags:
	ctags -R .

.PHONY: lint
lint:
	golongfuncs
	gometalinter --deadline=60s --disable=interfacer gpx

.PHONY: install-tools
install-tools:
	go get -u github.com/tkrajina/golongfuncs/...
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter --install