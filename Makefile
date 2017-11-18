GO=go
GOFLAGS=build
GOSRCS=wololo.go logging.go config.go networking.go
TESTSRCS=Dockerfile setup_cont_env.sh
BINNAME=wololo
IMGNAME=wololo-test-img
CONTNAME=wololo-test
NETNAME=wolnet

.PHONY : all
all : wololo test

# Target to build binary
$(BINNAME): $(GOSRC)
	$(GO) $(GOFLAGS) $(GOSRCS)

# Target to build test container
.PHONY: test
test: $(TESTSRCS) $(BINNAME)
	# Build docker container image
	sudo docker build --rm -t $(IMGNAME) .

	# Set up the test network
	sudo docker network create $(NETNAME)

	# Create the container and connect to network
	sudo docker create --name $(CONTNAME) $(IMGNAME)
	sudo docker network connect $(NETNAME) $(CONTNAME)

.PHONY: clean
clean:
	# Remove binary
	rm $(BINNAME)

	# Docker container related teardown
	sudo docker stop $(CONTNAME)
	sudo docker rm $(CONTNAME)
	sudo docker rmi $(IMGNAME)
	sudo docker network rm $(NETNAME)
