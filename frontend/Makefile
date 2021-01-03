build:
	npm run build:prod

copyToDist: build
	mkdir -p ../dist/public
	cp -r dist/monhttp/ ../dist/public/

clean:
	rm -rf dist
