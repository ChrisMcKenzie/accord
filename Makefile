PLATFORMS := darwin_amd64 linux_amd64

GIT_COMMIT:=$(shell git rev-parse HEAD)
GIT_DIRTY:=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

build/accord_$(PLATFORMS):
	gox $(foreach platform,$(subst _,/,$(PLATFORMS)),--osarch="$(platform)") \
		--output="build/accord_{{.OS}}_{{.Arch}}" \
		--ldflags="-X github.com/datascienceinc/accord/cmd.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X github.com/datascienceinc/accord/cmd.Version=$(VERSION) -X github.com/datascienceinc/accord/cmd.BuildDate=$(shell date -u +.%Y%m%d.%H%M%S)"		
	

build: build/accord_$(PLATFORMS)

clean:
	rm -rf build .accord
