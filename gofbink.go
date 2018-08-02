package gofbink

// #cgo LDFLAGS: -L${SRCDIR}/fbinklib -lfbink
// #include <stdlib.h>
// #include "fbinkinclude/fbink.h"
import "C"
import (
	"errors"
	"unsafe"
)

type FBInkConfig struct {
	row         int16
	col         int16
	fontmult    uint8
	fontname    uint8
	isInverted  bool
	isFlashing  bool
	isCleared   bool
	isCentered  bool
	isPadded    bool
	isVerbose   bool
	isQuiet     bool
	ignoreAlpha bool
}

func fbconfigGoToC(fbConf FBInkConfig) C.FBInkConfig {
	var cFBconfig C.FBInkConfig
	cFBconfig.row = C.short(fbConf.row)
	cFBconfig.col = C.short(fbConf.col)
	cFBconfig.fontmult = C.uint8_t(fbConf.fontmult)
	cFBconfig.fontname = C.uint8_t(fbConf.fontname)
	cFBconfig.is_inverted = C.bool(fbConf.isInverted)
	cFBconfig.is_flashing = C.bool(fbConf.isFlashing)
	cFBconfig.is_cleared = C.bool(fbConf.isCleared)
	cFBconfig.is_centered = C.bool(fbConf.isCentered)
	cFBconfig.is_padded = C.bool(fbConf.isPadded)
	cFBconfig.is_verbose = C.bool(fbConf.isVerbose)
	cFBconfig.is_quiet = C.bool(fbConf.isQuiet)
	cFBconfig.ignore_alpha = C.bool(fbConf.ignoreAlpha)
	return cFBconfig
}

func fbinkVersion() string {
	vers := C.GoString(C.fbink_version())
	return vers
}

func fbinkInit(fbfd int, cfg FBInkConfig) error {
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

func fbinkPrint(fbfd int, str string, cfg FBInkConfig) error {
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

func fbinkRefresh(fbfd int, top, left, width, height uint32, waveMode string, blackFlash bool) error {
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
func fbinkIsFBquirky() bool {
	var resultC C.bool
	resultC = C.fbink_is_fb_quirky()
	return bool(resultC)
}

func fbinkPrintImage(fbfd int, imgPath string, targX, targY int16, cfg FBInkConfig) error {
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
