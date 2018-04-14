build:
	go build main.go token.go widgets.go user.go helper.go

run:
	./main

clean:
	go clean

all:
	make clean
	make build
	make run

