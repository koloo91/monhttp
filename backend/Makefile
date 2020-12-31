build:
	go build -o monhttp

copyToDist: build
	mkdir -p ../dist
	cp monhttp ../dist/monhttp
	mkdir -p ../dist/migrations/
	cp -r migrations/ ../dist/migrations/

clean:
	rm -f monhttp
