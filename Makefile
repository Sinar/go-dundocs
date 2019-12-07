run: build
	@./dundocs

build:
	@go build -o dundocs ./cmd/go-dundocs 

# run: build
# 	@./go-pardebate
#
# build:
# 	@go build ./cmd/go-pardebate
#
# runpardocs: buildpardocs
# 	@./go-pardocs
#
# buildpardocs:
# 	@go build ./cmd/go-pardocs

test:
	@gotest ./...

# https://vic.demuzere.be/articles/golang-makefile-crosscompile/
PLATFORMS := linux/amd64 windows/amd64 darwin/amd64 freebsd/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	@mkdir -p dist/$(os)-$(arch)
	@GOOS=$(os) GOARCH=$(arch) go build -o './dist/$(os)-$(arch)/go-dundocs_$(os)-$(arch)' ./cmd/go-dundocs
	# @GOOS=$(os) GOARCH=$(arch) go build -o './dist/$(os)-$(arch)/go-pardebate_$(os)-$(arch)' ./cmd/go-pardebate

.PHONY  release: $(PLATFORMS)
