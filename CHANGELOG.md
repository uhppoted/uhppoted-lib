## CHANGELOG

### v0.7.2

1. Reworked `config` to support types.BindAddr, types.BroadcastAddr and types.ListenAddr
2. Implemented shared `config` command
3. Removed rollover from event handling (cf. https://github.com/uhppoted/uhppote-cli/issues/7)
4. Removed `EventRange` function as not practical.
5. Added `GetNextEvent` function to simplify sequential event retrieval

### v0.7.1

1. Implemented `PutTaskList` function
2. Reworked `UHPPOTED` as an interface+implementation
3. `encoding/conf` moved from `uhppote-core`

