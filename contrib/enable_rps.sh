#!/bin/sh
#=====================================
# Enable RPS (Receive Packet Steering)
#
# rps_cpus=3 because I have 2 cpu cores
#       
#       bin
# CPU0  01
# CPU1  10
# --------
#       11
#
#=====================================

rfc=4096
cc=$(grep -c processor /proc/cpuinfo)
rsfe=$(echo $cc*$rfc | bc)
sysctl -w net.core.rps_sock_flow_entries=$rsfe
for fileRps in $(ls /sys/class/net/eth*/queues/rx-*/rps_cpus)
do
    echo 3 > $fileRps
done
for fileRfc in $(ls /sys/class/net/eth*/queues/rx-*/rps_flow_cnt)
do
    echo $rfc > $fileRfc
done

tail /sys/class/net/eth*/queues/rx-*/{rps_cpus,rps_flow_cnt}

watch -d -n 1 'cat /proc/softirqs'

