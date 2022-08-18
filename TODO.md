# TODO

### IN PROGRESS

- [x] Once off crash in tests after updating to Go 1.19:
```
go test ./...
fatal error: concurrent map writes

goroutine 83 [running]:
github.com/uhppoted/uhppoted-lib/acl.TestPutACLWithConcurrency.func3(0x1881b0?, {0x0?, 0xc00018e258?, 0xc00018e270?, 0xc0001881b0?})
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put_test.go:564 +0xa8
github.com/uhppoted/uhppoted-lib/acl.(*mock).PutCard(0x11e70e0?, 0x1b8000?, {0x1?, 0xc00018e258?, 0xc00018e270?, 0xc0001881b0?})
```

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
