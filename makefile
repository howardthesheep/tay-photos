gofiles := $(shell cd server && find . -type f -name '*.go')

run:
	cd server && go build -o ../build/server $(gofiles)
	./build/server