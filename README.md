Wololo - Wake on Lan Service
======
**Synopsis** Device is off - WOLOLO - Device is on!

## Table of Contents
* Description
* Building and Testing
* Installation and Configuration
* TODOs

## Description
Wololo is a Wake on Lan (WOL) service application that runs as a simple webserver. The device it is intended to wake up, as well as the address and port it listens on can be configured. When the web service is accessed, it will send a WOL signal with the configured MAC address. This allows Wololo to either be run locally on your machine, or on another device in your network in order to wake the device up remotely.

## Building and Testing
Wololo can be built such that only the binary is output, or such that a test environment for debugging and development is also set up.

#### Requirements
The following required versions are based on those used when developing this project.
GNU Make version >= 4.1 is required when building using the provided Makefile.
To build the Wololo binary, a go version >= 1.6.3 is required.
To be able to build and run the test environment, Docker version >= 17.05.0-ce will be needed.


#### Building the Binary
Building go is as simple as invoking make with the target for the binary.

```bash
make wololo
```

#### Building the Test Environment
Building the test environment will build a Docker container which includes the Wololo application.
In addition, the container is attached to a newly created network which emulates the connection to the device to wake up.
As a result, the container will see two network interfaces. One for listening for requests by the used, and another for sending the WOL packet.

```bash
$ make wololo
$ make test
```

#### Running the Test

First, the Docker container built in the previous step has to be started. This will setup Wololo inside the container and extract relevant information from the container network. The network interface information from the container is printed to the user. The following shows an example.
```bash
$ sudo docker start -i wololo-test
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 172.17.0.2  netmask 255.255.0.0  broadcast 0.0.0.0
        ether 02:42:ac:11:00:02  txqueuelen 0  (Ethernet)
        RX packets 3  bytes 290 (290.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

eth1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 172.18.0.2  netmask 255.255.0.0  broadcast 0.0.0.0
        ether 02:42:ac:12:00:02  txqueuelen 0  (Ethernet)
        RX packets 2  bytes 180 (180.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        loop  txqueuelen 1  (Local Loopback)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

Note that Wololo will be listening on eth0 (i.e. 172.17.0.2), on pot 5000. This is the IP and port that will be used to trigger the service.
Before doing this, the unit test must be started. This will create a UDP listener that will wait for the WOL signal. Run the following in another console.

```bash
$ sudo go test wololo_test.go
```

In a third console, request the page served by the Wololo service in the Docker container.

```bash
$ curl 172.17.0.2:5000
Device is off...
WOLOLO
Device is on!
```

This indicates that the WOL sequence has been sent. Check the console of the unit test, whether the correct sequence was received. The following shows the output in case everything went fine.

```bash
$ sudo go test wololo_test.go
ok  	command-line-arguments	4.028s
```

To stop the Docker container, run the following command.


```bash
$ sudo docker stop wololo-test
```
## Installation and Configuration
Once built, the Wololo service must be configured before invoking the binary. The corresponding configuration file is expected under /etc/wololo/wololo.conf. It can be created as follows.

```bash
$ sudo mkdir /etc/wololo
$ sudo touch /etc/wololo/wololo.conf
$ sudo chmod a+r /etc/wololo/wololo.conf
```
The following shows an exemplary configuration file.

```bash
$ cat /etc/wololo/wololo.conf
Listen=172.17.0.2:5000
Broadcast=172.18.255.255:7
Interface=eth1
MAC=00:11:22:33:44:55:66
```

The 'Listen' parameter defines the IP and port on which the Wololo service will listen. The 'Broadcast' parameter defines the subnet and port to broadcast the WOL signal on. Typically the port will be 7 or 9. The 'Interface' specifies the physical interface on which the signal should be sent. The 'MAC' parameter is the MAC address ot the device to wake up.

The only strictly required parameter is the 'MAC' parameter. All other parameters will default to a hardcoded value. These are as follows.
* Listen: 127.0.0.1:5000
* Broadcast: 255.255.255.255:7
* Interface: eth0

## TODOs
The following list is a set of tasks that should be done at some point.
* Integrate an HTTP GET into the wololo_test.go unit test, rather than having to use curl to trigger the test
* Package and provide Wololo as a snap using snapcraft
