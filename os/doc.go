// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
The os package implements replacement functions for the Go [os/https://pkg.go.dev/os] package.

It currently only implemants a replacement function for os.Rename to workaround the
'invalid cross-device link' error when renaming across file-systems

	(cf. https://github.com/uhppoted/uhppoted-httpd/issues/20)
*/
package os
