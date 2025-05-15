NAME = "meermookh"

build:
	go build -o bin/$(NAME) 

run: build
	./bin/$(NAME)
