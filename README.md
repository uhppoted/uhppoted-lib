![build](https://img.shields.io/badge/dynamic/json?url=https://gist.github.com/twystd/2609d9356470acaf904386ce4a74f54e/raw/uhppoted-lib.json&query=$.message&label=build&color=brightgreen)


# uhppoted-lib

Shared library that implements the functionality common to multiple _uhppoted_ sub-projects. In particular:

- _conf_ file marshaling
- ACL (access control lists)
- TSV encoding and decoding
- Wrapper functions for the rest and MQTT services to facilitate invoking controller functions on multiple devices.

## Release Notes

### Current Release

**[v0.9.0](https://codeberg.org/uhppoted/uhppoted-lib/releases/tag/v0.9.0) - 2026-01-27**

1. Updated to Go 1.25.


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
