NAME=monhttp
VERSION=`cat version`

.PHONY: buildDockerImage
buildDockerImage:
	docker build -t koloooo/monhttp:$(VERSION) -t koloooo/monhttp:latest .

.PHONY: buildAndPushDockerImage
buildAndPushDockerImage:
	docker build -t koloooo/monhttp:$(VERSION) -t koloooo/monhttp:latest .
	docker push --all-tags koloooo/monhttp

.PHONY: incrementReleaseVersion
incrementReleaseVersion:
	./increment_version.sh -m $(VERSION) > version
	$(eval VERSION=`cat version`)
	echo $(VERSION)

.PHONY: release
release: incrementReleaseVersion buildAndPushDockerImage

.PHONY: buildLocal
buildLocal:
	rm -rf dist
	$(MAKE) -C backend copyToDist
	$(MAKE) -C frontend copyToDist
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean
