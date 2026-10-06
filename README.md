![build](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/twystd/2609d9356470acaf904386ce4a74f54e/raw/uhppoted-lib.json)


# uhppoted-lib

Shared library that implements the functionality common to multiple _uhppoted_ sub-projects. In particular:

- _conf_ file marshaling
- ACL (access control lists)
- TSV encoding and decoding
- Wrapper functions for the rest and MQTT services to facilitate invoking controller functions on multiple devices.

## Release Notes

### Current Release

## [0.9.1](https://codeberg.org/uhppoted/uhppoted-lib/releases/tag/v0.9.1) - 2026-10-06

1. Migrated from _github_ to _codeberg_.
2. Updated to Go v1.27.
3. Added _SetFirstCard_ API function.
4. Added _first card_ privileges to card.
5. (Optionally) included _PIN_ and _first-card_ privileges in _ACL_ APIs.
6. Added `rest.PIN.enabled` and  `rest.firstcard.enabled` to _uhppoted.conf_.
7. Added `mqtt.PIN.enabled` and  `mqtt.firstcard.enabled` to _uhppoted.conf_.


## Development

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://codeberg.org/uhppoted/uhppoted-lib.git
cd uhppoted-lib
make build
```

If you prefer not to use `make`:
```
git clone https://codeberg.org/uhppoted/uhppoted-lib.git
cd uhppoted-lib
mkdir bin
go build -trimpath -o bin ./...
```
