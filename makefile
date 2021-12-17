gofiles := $(shell find . -type f -name '*.go')

run: $(gofiles)
	go build -o build/server $(gofiles)
	cp -R www/. build/www/
	./build/server