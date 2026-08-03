BINARY = espresso
BUNDLE = Espresso.app

.PHONY: build clean

build:
	go build -o $(BINARY) .
	mkdir -p $(BUNDLE)/Contents/MacOS
	mkdir -p $(BUNDLE)/Contents/Resources
	cp $(BINARY) $(BUNDLE)/Contents/MacOS/
	cp Info.plist $(BUNDLE)/Contents/
	cp AppIcon.icns $(BUNDLE)/Contents/Resources/
	printf 'APPL????' > $(BUNDLE)/Contents/PkgInfo
	@echo "Built $(BUNDLE)"

clean:
	rm -rf $(BINARY) $(BUNDLE) AppIcon.iconset
