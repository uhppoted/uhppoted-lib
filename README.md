![build](https://github.com/uhppoted/uhppoted-lib/workflows/build/badge.svg)

# uhppoted-lib

Shared library that implements the functionality common to multiple _uhppoted_ sub-projects. In particular:

- _conf_ file marshaling
- ACL (access control lists)
- TSV encoding and decoding
- Wrapper functions for the rest and MQTT services to facilitate invoking controller functions on multiple devices.

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.8.5    | Added `SetInterlock` wrapper API for `set-interlock` function                             |
| v0.8.4    | Added support for card keypad PIN to ACL                                                  |
| v0.8.3    | Added lockfile implementation using `flock` _syscall_                                     |
| v0.8.2    | Fixed address resolution bug in health-check                                              |
| v0.8.1    | Maintenance release for version compatibility with `uhppote-core` `v0.8.1`                |
| v0.8.0    | Maintenance release for version compatibility with `uhppote-core` `v0.8.0`                |
| v0.7.3    | Maintenance release for version compatibility with `uhppote-core` `v0.7.3`                |
| v0.7.2    | Replaced event rollover with infinite event indexes to match controller implementation    |
| v0.7.1    | Added support for task list functions from the extended API                               |
| v0.7.0    | Added support for time profiles from the extended API                                     |
| v0.6.12   | Additional validation of bind, broadcast and listen ports when loading configuration      |
| v0.6.10   | Adds configuration options for initial release of `uhppoted-app-wild-apricot`             |
| v0.6.8    | Maintenance release for version compatibility with `uhppote-core` `v0.6.8`                |
| v0.6.7    | Implements `record-special-events` for enabling/disabling door events                     |
| v0.6.5    | Maintenance release for version compatibility with NodeRED module                         |
| v0.6.4    | Added support for uhppoted-app-sheets                                                     |
| v0.6.3    | Added support for `uhppoted-mqtt` ACL API                                                 |
| v0.6.2    | Added support for `uhppoted-rest` ACL API                                                 |
| v0.6.1    | Added support for `uhppote-cli` ACL functions                                             |
| v0.6.0    | Added support for `uhppoted-acl-s3` ACL functions                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

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




