.PHONY: all build clean fmt bootstrap lint test
ALL_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
ALL_PKGS := $(shell glide nv |xargs go list)
NAME := $(shell grep "^[^=]" README.rst |head -1)
ALL_COVER := $(addsuffix /cover.out,$(subst github.com/Robpol86/${NAME},.,${ALL_PKGS}))

all: clean lint test

clean:
	rm -f $(ALL_COVER)

$(GOPATH)/bin/golint:
	go get -u github.com/golang/lint/golint

lint: $(GOPATH)/bin/golint
	@echo "Running golint"
	echo $(ALL_PKGS) |xargs -n1 golint |(! grep --color '.')
	@echo "Running go vet"
	go vet $(ALL_PKGS)
	@echo "Checking gofmt"
	gofmt -l $(ALL_FILES) |(! grep --color '.')

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide

bootstrap: $(GOPATH)/bin/glide

${ALL_COVER}: PKG=$(addprefix github.com/Robpol86/${NAME}/,$(dir $@))
${ALL_COVER}:
	go test -coverprofile $@ $(PKG)

test: clean ${ALL_COVER}
	go version

fmt:
	@echo Formatting Packages...
	gofmt -l $(ALL_FILES) |xargs -L1 go fmt
