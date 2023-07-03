# CHANGELOG

## Unreleased

### Updated
1. Added _card.format_ to the configuration to facilitate support for card formats other
   than Wiegand-26.
1. Added _card formats_ parameter to load-acl for optional card number validation.
   

## [0.8.5](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.5) - 2023-06-13

### Added
1. `SetInterlock` API function for REST and MQTT `set-interlock` command.

### Updated
1. Reworked to use cards with pointerless 'from' and 'to' dates


## [0.8.4](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.4) - 2023-03-17

### Added
1. `doc.go` package overview documentation

### Updated
1. Fixed initial round of _staticcheck_ lint errors and permanently added _staticcheck_ to
   CI build.
2. Reworked ACL functions to retain stored card PIN

## [0.8.3](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Lockfile implementation based on the `flock` _syscall_ (`FileLock` and `FileUnlock` on _Windows_).
2. Added `mqtt.connection.verify` to configuration for _uhppoted-mqtt_.

### Changed
1. Added OTP section to _uhppoted-httpd_ config.
2. Suppressed '... displaying configuration' message from _config_ command.


## [0.8.2](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.1) - 2022-10-14

### Changed
1. Explicitly converted listener addresses to IPv4 addresses in health check (Ref. [Fix health-check listener monitor to accommodate the default 0.0.0.0 address](https://github.com/uhppoted/uhppoted-lib/issues/2#issuecomment-1204253581)).
2. Updated go.mod to Go 1.19.
3. Added softlock to MQTT configuration.
4. Added ACL mqtt.acl.verify to MQTT configuration.
5. Reworked RecordSpecialEvents to not use wrapped requests/responses.
6. Reworked watchdog to account for configurable healthcheck interval.
7. Included health-check interval in watchdog configuration. 

## [0.8.1](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Fixed panic in GetCard (dereferencing invalid card)
2. Simplified GetStatus API
3. Added locales and _en_ dictionary for event type and reason lookup.
4. Removed wrapper for events received on by `Listen`
5. Added support to load locale translation dictionary from JSON file.
6. Added protocol.version and translation.locale to REST and MQTT configurations.
7. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.

## [v0.8.0](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.0) - 2022-07-01

### Changed
1. Updated uhppote-core dependency to v0.8.0

## [0.7.3](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.7.3) - 2022-06-01

### Changed
1. Added `SetEventListener` function.
2. Updated for reworked types.Date and types.DateTime zero values

## [0.7.2](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.7.2)

### Changed
1. Reworked `config` to support types.BindAddr, types.BroadcastAddr and types.ListenAddr
2. Implemented shared `config` command
3. Removed rollover from event handling (cf. https://github.com/uhppoted/uhppote-cli/issues/7)
4. Removed `EventRange` function as not practical.
5. Added `GetNextEvent` function to simplify sequential event retrieval

## [0.7.1](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.7.1)

### Changed
1. Implemented `PutTaskList` function
2. Reworked `UHPPOTED` as an interface+implementation
3. `encoding/conf` moved from `uhppote-core`

