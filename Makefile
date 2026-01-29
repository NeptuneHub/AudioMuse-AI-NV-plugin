PLUGIN_NAME := audiomuseai
PACKAGE_NAME := audiomuseai

.PHONY: build package clean

build:
	tinygo build -o $(PLUGIN_NAME).wasm -target=wasip1 -scheduler=none -buildmode=c-shared main.go

package: build
	# Navidrome expects the wasm file in the package to be named `plugin.wasm`
	cp $(PLUGIN_NAME).wasm plugin.wasm
	zip -j $(PACKAGE_NAME).ndp manifest.json plugin.wasm
	# Clean up temporary and built wasm output
	rm -f plugin.wasm $(PLUGIN_NAME).wasm

clean:
	rm -f $(PLUGIN_NAME).wasm $(PACKAGE_NAME).ndp
