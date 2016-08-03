#!/bin/sh

for i in `seq 1 30`; do 
    sh -c 'echo -n `date` " " ; ifconfig eth0 ; sleep 1 ; ifconfig eth0' | grep bytes | sed -e 's/.. bytes://g' | awk 'NR == 1 { rx = $1; tx = $4} NR == 2 { printf "rx %.1f Mbit/sec tx %.1f Mbit/sec ", ($1 - rx) * 8 / (1024*1024) / 1, ($4 - tx) * 8 / (1024*1024) / 1 }' ; date
done
