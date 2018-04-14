build:
	go build main.go TokenHelper.go WidgetsHelper.go UserHelper.go helper.go

run:
	./main

clean:
	go clean

all:
	make clean
	make build
	make run

