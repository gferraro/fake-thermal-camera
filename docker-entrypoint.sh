
#!/bin/bash
echo --- device-register ----
cd /code/device-register
echo Building device-register ....
go build ./...

echo --dbus --
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



echo --- thermal-recorder ----
cd /code/thermal-recorder/cmd/thermal-recorder
echo Building thermal-recorder....
go build ./...
cp    "../../_release/org.cacophony.thermalrecorder.conf" "/etc/dbus-1/system.d/org.cacophony.thermalrecorder.conf"
echo Running thermal recorder...
# ./thermal-recorder &
# disown

echo --- thermal-uploader ----
cd /code/thermal-uploader/
echo Building thermal-uploader....
go build ./...
echo Running thermal uploader...
./thermal-uploader &
disown


echo --- test-server ----
cd /server
echo Building testing-server....
go build ./...
echo Running test server...
cd cmd/testing-server/
./testing-server