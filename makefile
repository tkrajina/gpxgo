build:
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
prepare:
	go get
clean:
	echo "TODO"
