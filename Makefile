test:
	./test.sh

build:
	mkdir ./dist
	go build -o universe
	mv ./universe ./dist/
	cp ./config.json ./dist/

clean:
	rm -rf ./dist

docker:
	GOOS=linux make build
	cp ./Dockerfile ./dist
	docker build --no-cache -t universe:latest ./dist
