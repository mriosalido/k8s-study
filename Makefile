
VERSION=1.2

build:
	sed -i -e 's/var Version = ".*/var Version = "$(VERSION)"/' version/version.go
	go build

image: build
	docker build -t mriosalido/k8s-study:$(VERSION) .

push: image
	docker push mriosalido/k8s-study:$(VERSION)
