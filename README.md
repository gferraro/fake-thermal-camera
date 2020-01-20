# fake thermal camera

`fake-thermal-camera` is a server that runs some of the thermal camera code (currently on linux not raspberian).   Is is designed for using with the cypress integration tests.

Project | fake-thermal-camera
---|---
Platform | Linux
Requires | Git repository [`cacophony-api`](https://github.com/TheCacophonyProject/cacophony-api) server to connect to </br> Git repository [`device-register`](https://github.com/TheCacophonyProject/device-register)
Licence | GNU General Public License v3.0

## Development Instructions

Download fake-thermal-camera, and device-register project into the same folder.

Make sure the [`cacophony-api`](https://github.com/TheCacophonyProject/cacophony-api) server the devices will attach to is also running

In the fake-thermal-camera folder start the test server with
```
> ./run
```

Now you can start calling server commands

## Current server commands
```
GET http://localhost:2040/create/{device-name}?group-name={group-name}  Create a new device.  Needs to be called before any other command.
```

Note that the server currently only runs one device at a time. Once a new device is created it is no longer possible to use old devices.
