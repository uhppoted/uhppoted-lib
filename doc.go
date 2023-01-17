// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppoted-lib implements the functionality common to most applications that use uhppote-core.

The package is a library package that is intended to simplify the common tasks associated with writing
applications that interface to the UHPPOTE access controllers but that do not form part of the core API,
viz.:
- configuration
- command line arguments
- access control lists
- internationalisation
- lock files
- logging
- monitoring

# Structure

The library consists of the following packages:
  - [uhppoted-lib/uhppoted] which wraps certain uhppote-core functionality that typically requires an
    extended interchange with the access controller e.g. retrieving card lists or events.
  - [uhppoted-lib/conf] which manages the common configuration file uhppoted.conf.
  - [uhppoted-lib/command] which defines the common 'help' and 'version' commands along with helper functions
    for command line parsing.
  - [uhppoted-lib/acl] which implements the commonly required access control list management functionality.
  - [uhppoted-lib/log] which implements the common logging format used by other uhppoted modules.
  - [uhppoted-lib/lockfile] which implements the lockfiles used to ensure single active instances of an application.
  - [uhppoted-lib/monitoring] which implements the system health and watchdog functionality.
*/
package lib
