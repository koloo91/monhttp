NAME=monhttp
VERSION=`cat version`

REGISTRY=docker.pkg.github.com

.PHONY: buildAndPushDockerImage
buildAndPushDockerImage:
	docker build -t $(REGISTRY)/koloo91/monhttp/monhttp:$(VERSION) -t $(REGISTRY)/koloo91/monhttp/monhttp:latest .
	docker push --all-tags $(REGISTRY)/koloo91/monhttp/monhttp

.PHONY: incrementReleaseVersion
incrementReleaseVersion:
	./increment_version.sh -m $(VERSION) > version
	$(eval VERSION=`cat version`)
	echo $(VERSION)

.PHONY: release
release: incrementReleaseVersion buildAndPushDockerImage

.PHONY: buildLocally
buildLocally:
	rm -rf dist
	$(MAKE) -C backend copyToDist
	$(MAKE) -C frontend copyToDist
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean
