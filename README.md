![build](https://github.com/uhppoted/uhppoted-lib/workflows/build/badge.svg)

# uhppoted-lib

Shared library that implements the functionality common to multiple _uhppoted_ sub-projects. In particular:

- _conf_ file marshaling
- ACL (access control lists)
- TSV encoding and decoding
- Wrapper functions for the rest and MQTT services to facilitate invoking controller functions on multiple devices.

## Release Notes

### Current Release

**[v0.8.9](https://github.com/uhppoted/uhppoted-lib/releases/tag/v0.8.9) - 2024-09-06**

1. Added _httpd.security.admin.role_ to config.
2. Added default card start/end dates to HTTPD configuration.
3. Added admin role to HTTPD configuration.
4. Updated to Go 1.23.


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




