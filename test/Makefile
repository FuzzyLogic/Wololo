BINNAME=wololo
BINDIR=$(GOPATH)/bin
TESTSRCS=Dockerfile setup_cont_env.sh

.PHONY: all
all: test

# Target to build test container
.PHONY: test
test: $(TESTSRCS) $(BINDIR)/$(BINNAME)
	# Run tests
	sudo BINPATH=$(BINDIR)/$(BINNAME) ./run_tests.sh

# Cleanup in case of failure to build
.PHONY: clean
clean:
	if [ -f $(BINNAME) ]; then \
		rm $(BINNAME); \
	fi
