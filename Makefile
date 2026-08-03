BINARY = espresso
BUNDLE = Espresso.app

.PHONY: build clean

build:
	go build -o $(BINARY) .
	mkdir -p $(BUNDLE)/Contents/MacOS
	cp $(BINARY) $(BUNDLE)/Contents/MacOS/
	cp Info.plist $(BUNDLE)/Contents/
	@echo "Built $(BUNDLE)"

clean:
	rm -rf $(BINARY) $(BUNDLE)
