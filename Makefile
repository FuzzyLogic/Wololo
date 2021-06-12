.PHONY: build
build:
	go build -o build/wololo cmd/wololo/main.go

.PHONY: install
install:
	install -m 755 build/wololo /usr/local/bin/wololo
	mkdir -p /etc/wololo
	install -m 644 configs/config.json /etc/wololo/config.json

.PHONY: uninstall
uninstall:
	rm /usr/local/bin/wololo
	rm -rf /etc/wololo

.PHONY: clean
clean:
	rm build/*
