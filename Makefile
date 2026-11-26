.PHONY: dev build clean

dev:
	air

build:
	go build -o myapp .

clean:
	rm -rf tmp myapp