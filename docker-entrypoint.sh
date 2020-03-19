#!/bin/bash
echo --- device-register ----
cd /code/device-register
echo Building device-register ....
go build ./...

echo --dbus --
apt-get update && apt-get install -y dbus

mkdir -p /var/run/dbus
dbus-daemon --config-file=/usr/share/dbus-1/system.conf --print-address

echo --- event-reporter ----
cd /code/event-reporter
cd cmd/event-reporter/
echo Building event-reporter ....
go build

# cp event-report /usr/bin/event-reporter
# cd ../..
# cp _release/event-reporter.server /etc/systemd/system/event-reporter.service
# systemctl start event-report.service
# cd ../..
cp ../../_release/org.cacophony.Events.conf /etc/dbus-1/system.d/org.cacophony.Events.conf
./event-reporter --interval 2s &
disown
echo --- test-server ----
cd /server
echo Building testing-server....
go build ./...
echo Running test server...
./testing-server