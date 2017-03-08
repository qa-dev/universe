test:
	./test.sh

build:
	mkdir ./dist
	go build -o universe
	mv ./universe ./dist/
	cp ./config.yaml ./dist/

clean:
	rm -rf ./dist
