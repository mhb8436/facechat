#!/bin/bash
echo "### Kill all imserver processes"
tgtpids=`ps -j | grep main | grep module  | grep -v grep | awk '{print $2}'`
for tgtpid in $tgtpids
do
        echo "The main.go to be killed: $tgtpid"
        kill -9 $tgtpid
done
echo "### Start server"
sleep 3

go run main.go -module api > api.log 2>&1 &
go run main.go -module message > message.log 2>&1 &
go run main.go -module connect > connect.log 2>&1 &
go run main.go -module sender > sender.log 2>&1 &