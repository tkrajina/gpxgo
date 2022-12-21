test:
	go test ./gpx
gofmt:
	gofmt -w ./gpx
goimports:
	goimports -w ./gpx
install:
	go install ./gpx
prepare:
	go get
clean:
	echo "TODO"
ctags:
	ctags -R .
lint:
	golongfuncs
	gometalinter --deadline=60s --disable=interfacer gpx

install-tools:
	go get -u github.com/tkrajina/golongfuncs/...
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter --install