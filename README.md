# fake thermal camera

`fake-thermal-camera` is a server that runs some of the thermal camera code (currently on linux not raspberian).   Is is designed for using with the cypress integration tests.

Project | fake-thermal-camera
---|---
Platform | Linux
Requires | Git repository [`cacophony-api`](https://github.com/TheCacophonyProject/cacophony-api) server to connect to </br> Git repository [`device-register`](https://github.com/TheCacophonyProject/device-register) </br> Git repository [`thermal-recorder`](https://github.com/TheCacophonyProject/thermal-recorder) </br> Git repository [`thermal-uploader`](https://github.com/TheCacophonyProject/thermal-uploader) </br> Git repository [`management-interface`](https://github.com/TheCacophonyProject/management-interface)</br> Git repository [`management-interface`](https://github.com/TheCacophonyProject/event-reporter)</br> Git repository [`event-reporter`](https://github.com/TheCacophonyProject/management-interface)</br> Git repository [`management-interface`](https://github.com/TheCacophonyProject/management-interface)
Licence | GNU General Public License v3.0

## Development Instructions

Download fake-thermal-camera, and device-register, thermal-recorder, thermal-uploader, event-reporter, management-interface project into the same folder.

Make sure the [`cacophony-api`](https://github.com/TheCacophonyProject/cacophony-api) server the devices will attach to is also running

In the fake-thermal-camera folder start the test server with
```
> ./run
```

Now you can start calling server commands.  You should create a device before you get started.

## Current server commands
```
GET http://localhost:2040/create/{device-name}?group-name={group-name}  Create a new device.  Needs to be called before any other command.
GET http://localhost:2040/playCPTVFile/?cptv-file={filename}  Plays the CPTV (including telemetry data) through the camera (does not require a device to be created first)
GET http://localhost:2040/triggerEvent/{type}  Triggers and event of the given event type.
```

Note that the server currently only runs one device at a time. Once a new device is created it is no longer possible to use the previously created device.