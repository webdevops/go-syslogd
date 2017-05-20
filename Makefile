SOURCE = $(wildcard *.go)
TAG ?= $(shell git describe --tags)
GOBUILD = go build -ldflags '-w'

ALL = \
	$(foreach arch,64 32,\
	$(foreach suffix,linux osx,\
		build/go-logpipe-$(arch)-$(suffix))) \
	$(foreach arch,arm arm64,\
		build/go-logpipe-$(arch)-linux)

all: test build

build: clean test $(ALL)

# cram is a python app, so 'easy_install/pip install cram' to run tests
test:
	echo "No tests"
	#cram tests/*.test

clean:
	rm -f $(ALL)

# os is determined as thus: if variable of suffix exists, it's taken, if not, then
# suffix itself is taken
osx = darwin
build/go-logpipe-64-%: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=amd64 $(GOBUILD) -o $@

build/go-logpipe-32-%: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=386 $(GOBUILD) -o $@

build/go-logpipe-arm-linux: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $@

build/go-logpipe-arm64-linux: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $@

release: build
	github-release release -u webdevops -r go-logpipe -t "$(TAG)" -n "$(TAG)" --description "$(TAG)"
	@for x in $(ALL); do \
		echo "Uploading $$x" && \
		github-release upload -u webdevops \
                              -r go-logpipe \
                              -t $(TAG) \
                              -f "$$x" \
                              -n "$$(basename $$x)"; \
	done
