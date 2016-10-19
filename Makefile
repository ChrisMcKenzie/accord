PLATFORMS := darwin_amd64 linux_amd64

build/accord_$(PLATFORMS):
	gox $(foreach platform,$(subst _,/,$(PLATFORMS)),--osarch="$(platform)") \
		--output="build/accord_{{.OS}}_{{.Arch}}"		
	

build: build/accord_$(PLATFORMS)

clean:
	rm -rf build .accord
