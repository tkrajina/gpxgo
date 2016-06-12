test:
	go test ./gpx
gofmt:
	gofmt -w ./gpx
goimports:
	goimports -w ./gpx
build-generics:
	 gengen generic/nullable.go string \
            | gofmt -r 'NullableGeneric -> NullableString' \
            | gofmt -r 'NewNullableGeneric -> NewNullableString' \
                    > gpx/nullable_string.go
	 gengen generic/nullable.go int \
            | gofmt -r 'NullableGeneric -> NullableInt' \
            | gofmt -r 'NewNullableGeneric -> NewNullableInt' \
                    > gpx/nullable_int.go
	 gengen generic/nullable.go float64 \
            | gofmt -r 'NullableGeneric -> NullableFloat64' \
            | gofmt -r 'NewNullableGeneric -> NewNullableFloat64' \
                    > gpx/nullable_float64.go
	 gengen generic/nullable.go time.Time \
            | gofmt -r 'NullableGeneric -> NullableTime' \
            | gofmt -r 'NewNullableGeneric -> NewNullableTime' \
                    > gpx/nullable_time.go
install:
	go install ./gpx
prepare:
	go get
clean:
	echo "TODO"
ctags:
	ctags -R .
vet:
	go tool vet --all -shadow=true . 2>&1 | grep -v "declaration of err shadows"
lint:
	golint ./... | grep -v "or be unexported" | grep -v "underscores in"
errcheck:
	errcheck ./...
gocyclo:
	-gocyclo -over 10 .
check: test gocyclo vet lint errcheck
	echo "OK"
install-tools:
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/golang/lint
	go get -u github.com/kisielk/errcheck
