#!/bin/bash
#====================================
# show NIC PacketsPerSecond on Linux
#====================================
 
INTERVAL="1"  # update interval in seconds
 
if [ -z "$1" ]; then
        echo
        echo usage: $0 [network-interface]
        echo
        echo e.g. $0 eth0
        echo
        echo shows packets-per-second
        exit
fi
 
IF=$1
 
while true
do
        R1=`cat /sys/class/net/$1/statistics/rx_packets`
        T1=`cat /sys/class/net/$1/statistics/tx_packets`
        sleep $INTERVAL
        R2=`cat /sys/class/net/$1/statistics/rx_packets`
        T2=`cat /sys/class/net/$1/statistics/tx_packets`
        TXPPS_N=`expr $T2 - $T1`
        TXPPS=`expr $T2 - $T1 | sed ':a;s/\B[0-9]\{3\}\>/,&/;ta'`
        RXPPS_N=`expr $R2 - $R1` 
        RXPPS=`expr $R2 - $R1 | sed ':a;s/\B[0-9]\{3\}\>/,&/;ta'`
        PPS=`expr $TXPPS_N + $RXPPS_N | sed ':a;s/\B[0-9]\{3\}\>/,&/;ta'`
        echo "TX $1: $TXPPS pkts/s RX $1: $RXPPS pkts/s PPS: $PPS"
done
