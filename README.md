Wololo - Wake on Lan Service

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/wololo)

======
## Synopsys

In a nutshell:
> Device is off - WOLOLO - Device is on

## Table of Contents

1. [Description](#description)
2. [Building, running, installing and testing](#build-install-run-test)
3. [TODOs](#todos)

## Description <a name="description"></a>

Wololo is a Wake on Lan (WOL) service application that runs as a simple webserver. The device it is intended to wake up, as well as the address and port it listens on can be configured. When the web service is accessed, it will send a WOL signal with the configured MAC address. This allows Wololo to either be run locally on your machine, or on another device in your network in order to wake the device up remotely.

## Building, installing, running and testing <a name="build-install-run-test"></a>

### Building Wololo

To build the `wololo` command, run the following from the repository root.
```
$ make build
```
The binary is then located in the build/ directory.

Wololo can also be built as a Snap. In order to do this, run the following from the root of the repository.
```
$ snapcraft
```

### Installing Wololo

To install Wololo, run the following.
```
$ sudo make install
```

This install `wololo` to /usr/local/bin.
Make sure to adapt the configuration under /etc/wololo/config.json to your needs before running the software (see next section). 

### Running Wololo

The following shows the program's options.
```
$ wololo -h
Usage of wololo:
  -config string
        Path to Wololo configuration file (default "/etc/wololo/config.json")
```

The configuration file has the following parameters:
* `listenAddr`: The IP which the Wololo service will listen
* `listenPort`: The port on which the Wololo service listens
* `udpBcastAddr`: The subnet and port to broadcast the WOL signal on. Typically the port will be 7 or 9.
* `iface`: The physical interface on which the signal should be sent
* `macAddr`: The MAC address ot the device to wake up

### Testing Wololo

A simple scenario is provided with this repository to test the Wololo functionality without any hardware.
The test/ directory contains a Docker Compose setup consisting of a `device` and `wololo` service.
The `device` service emulates a device that receives the WOL signal.
The test protocol is as follows:
* The `wololo` and `device` services are started up
* The `wololo` service waits to be triggered to send the WOL signal via an HTTP GET request
* The `device` service sends the GET request to the `wololo` service
* The `wololo` service sends a pre-configured WOL signal to the `device` service
* The `device` service receives the pre-configured WOL signal and compares it to an expected reference
* If the received signal matches the reference, the test passes

To run the test, perform the following steps from the test/ directory and check the the output resembles the following.
```
$ docker-compose up
Creating network "test_testnet" with driver "bridge"
Creating wololo-test_wololo ... done
Creating wololo-test_device ... done
Attaching to wololo-test_wololo, wololo-test_device
wololo-test_device | 2021/06/12 19:19:07 Starting UDP server on port 7
wololo-test_device | 2021/06/12 19:19:07 Triggering Wololo service
wololo-test_device | 2021/06/12 19:19:07 Waiting for data
wololo-test_device | 2021/06/12 19:19:07 Received data on UDP connection! Checking...
wololo-test_device | 2021/06/12 19:19:07 Success!
```

## TODOs <a name="todos"></a>

The following list is a set of tasks that should be done at some point.
* Wake arbitry devices, based on contacted endpoint (config file only specifies default device to wake up)
