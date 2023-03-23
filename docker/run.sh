#!/usr/bin/env bash
set -e

PREFIX_EXEC=""
echo Spawning $NUM_DEVICES devices

umask 0000
pids=()
for ((i=0;i<$NUM_DEVICES;i++)); do
    LD_PRELOAD=/usr/local/lib/faketime/libfaketimeMT.so.1 /iotivity-lite/port/linux/service "device-$i" > /tmp/$i.log 2>&1 &
    pids+=($!)
done

/plgd/client-application $@ > /tmp/client-application.log 2>&1 &
pids+=($!)

terminate()
{
    echo "Terminate"
    for (( i=0; i<${#pids[@]}; i++ )); do
        kill -SIGTERM ${pids[$i]}
    done
}

trap terminate SIGTERM

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds
while sleep 10; do
for (( i=0; i<${#pids[@]}; i++ ));
do
    if ! kill -0 ${pids[$i]} 2>/dev/null; then
        echo "service[$i] with pid=${pids[$i]} is dead"
        exit 1
    fi
done
echo checking running devices
done