# timeular-linux

Simple Linux commandline client for the [Timeular](https://timeular.com/product/tracker/) tracking device.
Forked from [timeular-zei-linux](https://github.com/krisbuist/timeular-zei-linux).

Features:

- Usable locally without Timeular Web API, logs to filesystem
- Bluetooth LE connection via BlueZ DBUS interface
- Desktop notifications via DBUS

## Quickstart

1. Run `go get -u github.com/cschomburg/timeular-linux/cmd/timeular`
2. Grab the `config.json.example` and move to `~/.config/timeular/config.json`
3. Start the application: `./timeular`
