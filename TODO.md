# TODO

- [ ] TCP/IP protocol (cf. https://github.com/uhppoted/uhppote-core/issues/17)
      - [x] Add protocol to _conf_ file 'Devices' section
      - [ ] CHANGELOG
      - [ ] README

- [ ] Rework events handling (cf. https://github.com/uhppoted/uhppoted-mqtt/issues/16)
      - [ ] Remove unnecessary Event struct in uhppoted.events
            - (?) JSON tags to core type
      - [ ] Remove ListenBatchSize

## TODO

- [ ] config.NewConfig should not return pointer
- [ ] Simplify all the IUHPPOTED functions - no need to be so unsuccessfully generic
- [ ] Replace UHPPOTE parameter from ACL API with IUHPPOTED
- [ ] Rework PutTimeProfiles to return (response,BadRequestError) or somesuch rather than status code
- [ ] Rework Config to use plugins
      https://pkg.go.dev/plugin

1. Rework healthcheck to remove need for IUHPPOTE::DeviceList
2. Rework healthcheck to remove need for IUHPPOTE::ListenAddr
3. GetDevices: rename DeviceSummary.Address to IpAddress and use Address for IP+Port

### uhppoted-lib

- [ ] Logging
      - log.Warnf+ should default to stderr
      - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
      - Windows: event logging
      - Windows eventlog message file
        - https://social.msdn.microsoft.com/Forums/windowsdesktop/en-US/deaa0055-7770-4e55-a5b8-6d08b80b74af/creating-event-log-message-files
        - FormatMessage (https://go.dev/src/syscall/syscall_windows.go)

- [ ] websocket + GraphQL (?)
- [ ] IFTTT
- [ ] Braid (?)
- [ ] MacOS launchd socket handoff
- [ ] Linux systemd socket handoff
- [ ] conf file decoder: JSON
- [ ] Rework plist encoder
- [ ] move ACL and events to separate API's
- [ ] Make events consistent across everything
- [ ] Rework uhppoted-xxx Run, etc to use [method expressions](https://talks.golang.org/2012/10things.slide#9)
- [ ] system API (for health-check, watchdog, configuration, etc)
- [ ] Parallel-ize health-check 

### Documentation

- [ ] godoc
- [ ] build documentation

### Other

1. github project page
2. Integration tests
3. EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
4. TLA+/Alloy models:
    - watchdog/health-check
    - concurrent connections
    - HOTP counter update
    - key-value stores
    - event buffer logic
5. Update file watchers to fsnotify when that is merged into the standard library (1.4 ?)
    - https://github.com/golang/go/issues/4068
6. go-fuzz
