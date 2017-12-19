#!/bin/bash
BINNAME=wololo
IMGNAME=wololo-test-img
CONTNAME=wololo-test
NETNAME=wolnet

# Copy binary data for testing
cp $BINPATH .

# Copy libseccomp in case it was used
if [ "ARM" == "$GOARCH" ]; then
    if [ -f /usr/lib/arm-linux-gnueabihf/libseccomp.so.2.1.1 ]; then
        cp /usr/lib/arm-linux-gnueabihf/libseccomp.so.2.1.1 libseccomp.so.2
    fi
else
    if [ -f /lib/x86_64-linux-gnu/libseccomp.so.2.2.3 ]; then
        cp /lib/x86_64-linux-gnu/libseccomp.so.2.2.3 libseccomp.so.2
    fi
fi

# Build docker container image
docker build --rm -t $IMGNAME .

# Set up the test network
docker network create $NETNAME

# Create the container and connect to network
docker create --privileged --name $CONTNAME $IMGNAME
docker network connect $NETNAME $CONTNAME

# Run the container
docker start $CONTNAME

# Start the test
go test wololo_test.go

# Cleanup work
rm $BINNAME
rm libseccomp.so.2
docker stop $CONTNAME
docker rm $CONTNAME
docker rmi $IMGNAME
docker network rm $NETNAME
