## CHANGELOG

## [v0.8.1] - 2022-08-01

### Changed
1. Fixed panic in GetCard (dereferencing invalid card)
2. Simplified GetStatus API
3. Added locales and _en_ dictionary for event type and reason lookup.
4. Removed wrapper for events received on by `Listen`
5. Added support to load locale translation dictionary from JSON file.
6. Added protocol.version and translation.locale to REST and MQTT configurations.
7. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.

## [v0.8.0] - 2022-07-01

### Changed
1. Updated uhppote-core dependency to v0.8.0

## [v0.7.3] - 2022-06-01

### Changed
1. Added `SetEventListener` function.
2. Updated for reworked types.Date and types.DateTime zero values

## [v0.7.2]

### Changed
1. Reworked `config` to support types.BindAddr, types.BroadcastAddr and types.ListenAddr
2. Implemented shared `config` command
3. Removed rollover from event handling (cf. https://github.com/uhppoted/uhppote-cli/issues/7)
4. Removed `EventRange` function as not practical.
5. Added `GetNextEvent` function to simplify sequential event retrieval

## [v0.7.1]

### Changed
1. Implemented `PutTaskList` function
2. Reworked `UHPPOTED` as an interface+implementation
3. `encoding/conf` moved from `uhppote-core`

