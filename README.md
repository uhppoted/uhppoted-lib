![build](https://github.com/uhppoted/uhppoted-lib/workflows/build/badge.svg)

# uhppoted-lib

Shared library that implements the functionality common to multiple _uhppoted_ sub-projects. In particular:

- _conf_ file marshaling
- ACL (access control lists)
- TSV encoding and decoding
- Wrapper functions for the rest and MQTT services to facilitate invoking controller functions on multiple devices.

## Release Notes

### Current Release

**[v0.8.10](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.10) - 2025-01-29**

1. Added support for _auto-send interval_.


## Development

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppoted-lib.git
cd uhppoted-lib
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppoted-lib.git
cd uhppoted-lib
mkdir bin
go build -trimpath -o bin ./...
```

#### Dependencies

| *Dependency*                                             | *Description*                                          |
| -------------------------------------------------------- | ------------------------------------------------------ |
| [uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation                        |
| golang.org/x/sys                                         | Support for Windows services                           |




