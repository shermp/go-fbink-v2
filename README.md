# go-fbink-v2
go-fbink is a Go wrapper for the fbink tool found at https://github.com/NiLuJe/FBInk

go-fbink-v2 is no longer quite the simple wrapper anymore. It adds some convenience functions to use FBInk more like printing to a console. go-fbink-v2 now sets up an object with appropriate methods for use in programs.

go-fbink-v2 is currently tied to FBInk 1.10.3

## Installation and usage
go-fbink can be installed by doing the following:
```
go get github.com/shermp/go-fbink-v2
```
A precompiled static library is included for convenience if you are using a Kobo device. If you wish or need to compile your own library, you will need to make the "pic" target when compiling fbink from the FBInk directory

The static library should reside in `fbinklib/libfbink.a`

From your Go project, import go-fbink as follows:
```
import "github.com/shermp/go-fbink-v2/gofbink"
```
Note, you will need to enable cgo support when building your project, by setting the `CGO_ENABLED=1` environment variable when building, along with setting the `CC` and `CXX` environment variables to your ARM toolchain's GCC and G++ paths respectively.

A simple example program has been provided in `example/main.go`

You can refer to the original documentation found in the `fbink.h` file, which can be found at `fbinkinclude/fbink.h`.

The primary usage difference from FBInk is that where appropriate, go-fbink returns an error, or nil, rather than an integer to indicate success or failure. Note that the error string contains the C error code name (eg: "EXIT_FAILURE").

The only function that is unavailable in go-fbink is `fbink_printf()`. This is because cgo does not support variadic parameters. So, if that functionality is required, a simple `s := fmt.Sprintf("String %d", 1)` should do the trick...


