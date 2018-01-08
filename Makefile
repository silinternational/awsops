builddist:
	gox -output="dist/{{.OS}}/{{.Arch}}/awsops"

test:
	go test -cover
