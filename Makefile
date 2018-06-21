DOCKER_IMAGE = nagios-dnsblklist

all: buildimage
	docker run --rm -v $(CURDIR)/docker_build:/go/build $(DOCKER_IMAGE) \
		gox -osarch="linux/amd64 darwin/amd64" -ldflags="-s -w" \
  	-output="/go/build/{{.Dir}}_{{.OS}}_{{.Arch}}"

buildimage:
	docker build -t $(DOCKER_IMAGE) .
