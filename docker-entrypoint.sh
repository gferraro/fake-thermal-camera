#!/bin/bash
go version
echo --- device-register ----
cd /code/device-register
echo Building device-register ....
go build ./...

echo --starting dbus --
dbus-daemon --config-file=/usr/share/dbus-1/system.conf --print-address

echo --- event-reporter ----
cd /code/event-reporter
cd cmd/event-reporter/
echo Building event-reporter ....
go build
cp ../../_release/org.cacophony.Events.conf /etc/dbus-1/system.d/org.cacophony.Events.conf

echo --- thermal-recorder ----
cd /code/thermal-recorder/cmd/thermal-recorder
echo Building thermal-recorder....
go build ./...
cp    "../../_release/org.cacophony.thermalrecorder.conf" "/etc/dbus-1/system.d/org.cacophony.thermalrecorder.conf"
echo Running thermal recorder...

echo --- thermal-uploader ----
cd /code/thermal-uploader/
echo Building thermal-uploader....
go build ./...

echo --- fake-lepton ----
cd /server
echo Building fake-lepton....
cd cmd/fake-lepton/
go build
cp org.cacophony.FakeLepton.conf /etc/dbus-1/system.d/org.cacophony.FakeLepton.conf


echo --- management interface ----
cd /code/management-interface/
echo Building management-interface....
make build


echo --- starting supervisord ---
/usr/bin/supervisord &
disown

echo --- test-server ----
cd /server/cmd/testing-server/
echo Building test-server....
go build
echo Running test server...
./testing-server
