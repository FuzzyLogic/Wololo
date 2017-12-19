#!/bin/bash

# Configure the application
LISTENIP=$(ifconfig | grep -A1 eth0 | grep -oE '172.[0-9]+.[0-9]+.[0-9]')
BCASTIP=$(ifconfig | grep -A1 eth1 | grep -oE '172.[0-9]+')
echo "Listen=$LISTENIP:5000" >> /etc/wololo/wololo.conf
echo "Broadcast=$BCASTIP.255.255:7" >> /etc/wololo/wololo.conf
echo "Interface=eth1" >> /etc/wololo/wololo.conf
echo "MAC=00:11:22:33:44:55:66" >> /etc/wololo/wololo.conf

# Start application
#su testusr
ifconfig
strace /testdir/wololo
