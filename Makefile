builddist:
	gox -arch="amd64" -os="linux windows" -output="dist/{{.OS}}/{{.Arch}}/awsops"

test:
	go test -cover
