package util

import "runtime"

// RuntimeVersion Indicates the version of Go runtime the application is running on.
var RuntimeVersion = runtime.Version()
