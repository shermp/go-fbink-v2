package gofbink

// #cgo LDFLAGS: -L${SRCDIR}/fbinklib -lfbink
// #include <stdlib.h>
// #include "fbinkinclude/fbink.h"
import "C"
import (
	"errors"
	"unsafe"
)

// FBInkConfig is a struct which configures the behavior of fbink
type FBInkConfig struct {
	Row         int16
	Col         int16
	Fontmult    uint8
	Fontname    uint8
	IsInverted  bool
	IsFlashing  bool
	IsCleared   bool
	IsCentered  bool
	IsPadded    bool
	IsVerbose   bool
	IsQuiet     bool
	IgnoreAlpha bool
}

// fbconfigGoToC is a convenience function to convert our Go config struct
// to a C struct that fbink understands
func fbconfigGoToC(fbConf FBInkConfig) C.FBInkConfig {
	var cFBconfig C.FBInkConfig
	cFBconfig.row = C.short(fbConf.Row)
	cFBconfig.col = C.short(fbConf.Col)
	cFBconfig.fontmult = C.uint8_t(fbConf.Fontmult)
	cFBconfig.fontname = C.uint8_t(fbConf.Fontname)
	cFBconfig.is_inverted = C.bool(fbConf.IsInverted)
	cFBconfig.is_flashing = C.bool(fbConf.IsFlashing)
	cFBconfig.is_cleared = C.bool(fbConf.IsCleared)
	cFBconfig.is_centered = C.bool(fbConf.IsCentered)
	cFBconfig.is_padded = C.bool(fbConf.IsPadded)
	cFBconfig.is_verbose = C.bool(fbConf.IsVerbose)
	cFBconfig.is_quiet = C.bool(fbConf.IsQuiet)
	cFBconfig.ignore_alpha = C.bool(fbConf.IgnoreAlpha)
	return cFBconfig
}

// FBinkVersion gets the fbink version
func FBinkVersion() string {
	vers := C.GoString(C.fbink_version())
	return vers
}

// FBinkInit initializes the fbink global variables
// See "fbink.h" for detailed usage and explanation
func FBinkInit(fbfd int, cfg FBInkConfig) error {
	fbConf := fbconfigGoToC(cfg)
	fdC := C.int(fbfd)
	var resultC C.int
	resultC = C.fbink_init(fdC, &fbConf)
	res := int(resultC)
	if res < 0 {
		return errors.New("c function fbink_init encountered an error")
	}
	return nil
}

// FBinkPrint prints a string to the screen
// See "fbink.h" for detailed usage and explanation
func FBinkPrint(fbfd int, str string, cfg FBInkConfig) error {
	fbConf := fbconfigGoToC(cfg)
	fdC := C.int(fbfd)
	strC := C.CString(str)
	defer C.free(unsafe.Pointer(strC))
	var resultC C.int
	resultC = C.fbink_print(fdC, strC, &fbConf)
	res := int(resultC)
	if res < 0 {
		return errors.New("c function fbink_print encountered an error")
	}
	return nil
}

// FBinkRefresh provides a way of refreshing the eink screen
// See "fbink.h" for detailed usage and explanation
func FBinkRefresh(fbfd int, top, left, width, height uint32, waveMode string, blackFlash bool) error {
	fdC := C.int(fbfd)
	topC := C.uint32_t(top)
	leftC := C.uint32_t(left)
	widthC := C.uint32_t(width)
	heightC := C.uint32_t(height)
	waveModeC := C.CString(waveMode)
	defer C.free(unsafe.Pointer(waveModeC))
	blackFlashC := C.bool(blackFlash)
	var resultC C.int
	resultC = C.fbink_refresh(fdC, topC, leftC, widthC, heightC, waveModeC, blackFlashC)
	res := int(resultC)
	if res < 0 {
		return errors.New("c function fbink_refresh encountered an error")
	}
	return nil
}

// FBinkIsFBquirky tests for a quirky framebuffer state
// See "fbink.h" for detailed usage and explanation
func FBinkIsFBquirky() bool {
	var resultC C.bool
	resultC = C.fbink_is_fb_quirky()
	return bool(resultC)
}

// FBinkPrintImage will print an image to the screen
// See "fbink.h" for detailed usage and explanation
func FBinkPrintImage(fbfd int, imgPath string, targX, targY int16, cfg FBInkConfig) error {
	fdC := C.int(fbfd)
	imgPathC := C.CString(imgPath)
	defer C.free(unsafe.Pointer(imgPathC))
	xC := C.short(targX)
	yC := C.short(targY)
	fbConf := fbconfigGoToC(cfg)
	var resultC C.int
	resultC = C.fbink_print_image(fdC, imgPathC, xC, yC, &fbConf)
	res := int(resultC)
	if res < 0 {
		return errors.New("c function fbink_print_image encountered an error")
	}
	return nil
}
