build:
	go build main.go token.go widgets.go user.go helper.go

run:
	./main

clean:
	go clean

test:
	go test -cover

all:
	make clean
	make build
	make test
	make run

