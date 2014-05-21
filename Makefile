build:

	go build -o mani main.go

dist: build
	
	$(eval VER := $(shell ./mani -version))
	$(eval DISTPATH := dist/$(VER))

	gox -osarch="darwin/amd64 linux/amd64 windows/amd64"

	#
	# Creating Archive for $(VER)
	#
	
	mkdir -p $(DISTPATH)
	rm -rf $(DISTPATH)/*

	cp mani_darwin_amd64 mani_linux_amd64 mani_windows_amd64.exe $(DISTPATH)/;

.PHONY: build