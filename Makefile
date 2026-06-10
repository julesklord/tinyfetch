.PHONY: build install clean

BINARY=mini-fetch

build:
	cd $(CURDIR) && go build -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)

