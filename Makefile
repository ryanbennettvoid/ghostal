
clean:
	rm ./gho || true

build:
	go build -o gho ./cmd/main.go

install: clean build
	mv ./gho ~/go/bin