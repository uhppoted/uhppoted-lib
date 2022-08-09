# TODO

### IN PROGRESS

- [ ] Once off crash in tests after updating to Go 1.19:
```
go test ./...
fatal error: concurrent map writes

goroutine 83 [running]:
github.com/uhppoted/uhppoted-lib/acl.TestPutACLWithConcurrency.func3(0x1881b0?, {0x0?, 0xc00018e258?, 0xc00018e270?, 0xc0001881b0?})
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put_test.go:564 +0xa8
github.com/uhppoted/uhppoted-lib/acl.(*mock).PutCard(0x11e70e0?, 0x1b8000?, {0x1?, 0xc00018e258?, 0xc00018e270?, 0xc0001881b0?})
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/acl_test.go:78 +0x3a
github.com/uhppoted/uhppoted-lib/acl.putACL({0x11e70e0, 0xc0001b8000}, 0x107b437?, 0xc000047fb0?)
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put.go:113 +0x71f
github.com/uhppoted/uhppoted-lib/acl.PutACL.func1()
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put.go:42 +0x72
created by github.com/uhppoted/uhppoted-lib/acl.PutACL
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put.go:35 +0x219

goroutine 1 [chan receive]:
testing.(*T).Run(0xc000144000, {0x11aeaf5?, 0x5618e0ceeb4?}, 0x11b7a18)
    /usr/local/go/src/testing/testing.go:1494 +0x37a
testing.runTests.func1(0xc000144000?)
    /usr/local/go/src/testing/testing.go:1846 +0x6e
testing.tRunner(0xc000144000, 0xc000131cd8)
    /usr/local/go/src/testing/testing.go:1446 +0x10b
testing.runTests(0xc000000640?, {0x12bed40, 0x46, 0x46}, {0x133e108?, 0x40?, 0x12c1d40?})
    /usr/local/go/src/testing/testing.go:1844 +0x456
testing.(*M).Run(0xc000000640)
    /usr/local/go/src/testing/testing.go:1726 +0x5d9
main.main()
    _testmain.go:185 +0x1aa

goroutine 82 [semacquire]:
sync.runtime_Semacquire(0xc0001be050?)
    /usr/local/go/src/runtime/sema.go:62 +0x25
sync.(*WaitGroup).Wait(0xc000059a48?)
    /usr/local/go/src/sync/waitgroup.go:139 +0x52
github.com/uhppoted/uhppoted-lib/acl.PutACL({0x11e70e0?, 0xc0001b8000}, 0xd431?, 0x0)
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put.go:59 +0x354
github.com/uhppoted/uhppoted-lib/acl.TestPutACLWithConcurrency(0xc0001b2000)
    /Users/tonyseebregts/Development/uhppote/uhppoted/uhppoted-lib/acl/put_test.go:583 +0x21e9
testing.tRunner(0xc0001b2000, 0x11b7a18)
    /usr/local/go/src/testing/testing.go:1446 +0x10b
created by testing.(*T).Run
    /usr/local/go/src/testing/testing.go:1493 +0x35f
FAIL    github.com/uhppoted/uhppoted-lib/acl    1.652s
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
