PROJECT = $(shell basename $(CURDIR))

run:
	go build -o build/$(PROJECT) && ./build/$(PROJECT)
