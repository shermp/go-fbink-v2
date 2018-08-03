# go-fbink
go-fbink is a Go wrapper for the fbink tool found at https://github.com/NiLuJe/FBInk

## Installation and usage
go-fbink can be installed by doing the following:
```
go get github.com/shermp/go-fbink
// download the FBInk submodule
cd $GOPATH/src/github.com/shermp/go-fbink
git submodule update --init --recursive
```
A precompiled static library is included for convenience if you are using a Kobo device. If you wish or need to compile your own library, you will need to make the "pic" target when compiling fbink from the FBInk directory

The static library should reside in `fbinklib/libfbink.a`

From your Go project, import go-fbink as follows:
```
import gofbink "github.com/shermp/go-fbink"
```
Note, you will need to enable cgo support when building your project, by setting the `CGO_ENABLED=1` environment variable when building, along with setting the `CC` and `CXX` environment variables to your ARM toolchain's GCC and G++ paths respectively.

A simple example of usage is:
```
fbinkOpts := gofbink.FBInkConfig{4, 0, 0, 0, false, false, false, true, false, false, false, false}
gofbink.FBinkInit(-1, fbinkOpts)
gofbink.FBinkPrint(-1, "This is a test", fbinkOpts)

fbinkOpts.Row = 8
gofbink.FBinkPrint(-1, "This is another test", fbinkOpts)

gofbink.FBinkPrintImage(-1, "path/to/img.png", 10, 20, fbinkOpts)
```
You can refer to the original documentation found in the `fbink.h` file, which can be found at `fbinkinclude/fbink.h`. The usage is almost identical.

The primary usage difference is that where appropriate, go-fbink returns an error, or nil, rather than an integer to indicate success or failure.

The only function that is unavailable in go-fbink is `fbink_printf()`. This is because cgo does not support variadic parameters.
