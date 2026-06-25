BINARY=tinyfetch

.PHONY: build install clean test

# Build: prefer Go build if cmd/tinyfetch/main.go exists, otherwise no-op
build:
	@if [ -f "cmd/tinyfetch/main.go" ]; then \
		go build -o $(BINARY) ./cmd/tinyfetch; \
		echo "Built $(BINARY) (Go version)"; \
	else \
		echo "No Go implementation found — nothing to build"; \
	fi

# Install: install built binary
install: build
	@if [ -d "ascii" ]; then \
		mkdir -p /usr/local/share/tinyfetch/ascii; \
		cp -r ascii/* /usr/local/share/tinyfetch/ascii/; \
		echo "Installed ASCII assets to /usr/local/share/tinyfetch/ascii/"; \
	fi
	@if [ -f "$(BINARY)" ]; then \
		install -m 0755 $(BINARY) /usr/local/bin/$(BINARY); \
		echo "Installed built $(BINARY) binary to /usr/local/bin/$(BINARY)"; \
	else \
		echo "Nothing to install (binary not found)"; exit 1; \
	fi


# Test: runs the shell/Go test harness
test: build
	./tests/test.sh

clean:
	rm -f $(BINARY)


