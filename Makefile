all: go-get build

go-get:
	go get github.com/hashicorp/serf@v0.10.1
	go get github.com/twmb/franz-go/pkg/kgo@v1.13.5

build:
	mkdir -p bin/
	go build -o bin/squealy squealy.go

clean:
	rm bin/*
