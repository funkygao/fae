#!/bin/sh
#=====================================
# Enable RPS (Receive Packet Steering)
#=====================================

rfc=4096
cc=$(grep -c processor /proc/cpuinfo)
rsfe=$(echo $cc*$rfc | bc)
sysctl -w net.core.rps_sock_flow_entries=$rsfe
for fileRps in $(ls /sys/class/net/eth*/queues/rx-*/rps_cpus)
do
    echo fff > $fileRps
done
for fileRfc in $(ls /sys/class/net/eth*/queues/rx-*/rps_flow_cnt)
do
    echo $rfc > $fileRfc
done

tail -f /sys/class/net/eth*/queues/rx-*/{rps_cpus,rps_flow_cnt}
