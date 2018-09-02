/*
	FBInk: FrameBuffer eInker, a tool to print text & images on eInk devices (Kobo/Kindle)
	Copyright (C) 2018 NiLuJe <ninuje@gmail.com>

	go-fbink: A Go wrapper for FBInk
	Copyright (C) 2018 Sherman Perry

	----

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as
	published by the Free Software Foundation, either version 3 of the
	License, or (at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package gofbink

// #cgo LDFLAGS: -L${SRCDIR}/fbinklib -lfbink
// #include <stdlib.h>
// #include <errno.h>
// #include "FBInk/fbink.h"
import "C"
import (
	"errors"
	"unsafe"
)

// Font type
type Font uint8

// Font constants
const (
	IBM Font = iota
	UNSCII
	UNSCIIalt
	UNSCIIthin
	UNSCIIfantasy
	UNSCIImcr
	UNSCIItall
	Block
	Leggie
	Veggie
	Kates
	Fkp
	Ctrld
	Orp
	OrpB
	OrpI
	Scientifica
	ScientificaB
	ScientificaI
)

// Align type
type Align uint8

// Align index constants
const (
	None Align = iota
	Center
	Edge
)

// FGcolor type
type FGcolor uint8

// FGcolor constants
const (
	FGblack FGcolor = iota
	FGgray1
	FGgray2
	FGgray3
	FGgray4
	FGgray5
	FGgray6
	FGgray7
	FGgray8
	FGgray9
	FGgrayA
	FGgrayB
	FGgrayC
	FGgrayD
	FGgrayE
	FGwhite
)

// BGcolor type
type BGcolor uint8

// BGcolor constants
const (
	BGwhite BGcolor = iota
	BGgrayE
	BGgrayD
	BGgrayC
	BGgrayB
	BGgrayA
	BGgray9
	BGgray8
	BGgray7
	BGgray6
	BGgray5
	BGgray4
	BGgray3
	BGgray2
	BGgray1
	BGblack
)

// CexitCode type
type CexitCode int

// Go translation of FBInk's exit codes
const (
	exitSuccess = CexitCode(C.EXIT_SUCCESS)
	exitFailure = CexitCode(C.EXIT_FAILURE) * -1
	eNoDev      = CexitCode(C.ENODEV) * -1
	eNotSup     = CexitCode(C.ENOTSUP) * -1
)

// FBFDauto is the automatic fbfd handler
const FBFDauto = int(C.FBFD_AUTO)

// const exitSuccess = int(C.EXIT_SUCCESS)

// FBInkConfig is a struct which configures the behavior of fbink
type FBInkConfig struct {
	Row         int16
	Col         int16
	Fontmult    uint8
	Fontname    Font
	IsInverted  bool
	IsFlashing  bool
	IsCleared   bool
	IsCentered  bool
	Hoffset     int16
	Voffset     int16
	IsHalfway   bool
	IsPadded    bool
	FGcolor     FGcolor
	BGcolor     BGcolor
	IsOverlay   bool
	IsVerbose   bool
	IsQuiet     bool
	IgnoreAlpha bool
	Halign      Align
	Valign      Align
}

func createError(retValue CexitCode) error {
	switch retValue {
	case exitFailure:
		return errors.New("EXIT_FAILURE")
	case eNoDev:
		return errors.New("ENODEV")
	case eNotSup:
		return errors.New("ENOTSUP")
	default:
		return nil
	}
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
	cFBconfig.hoffset = C.short(fbConf.Hoffset)
	cFBconfig.voffset = C.short(fbConf.Voffset)
	cFBconfig.is_halfway = C.bool(fbConf.IsHalfway)
	cFBconfig.is_padded = C.bool(fbConf.IsPadded)
	cFBconfig.fg_color = C.uint8_t(fbConf.FGcolor)
	cFBconfig.bg_color = C.uint8_t(fbConf.BGcolor)
	cFBconfig.is_overlay = C.bool(fbConf.IsOverlay)
	cFBconfig.is_verbose = C.bool(fbConf.IsVerbose)
	cFBconfig.is_quiet = C.bool(fbConf.IsQuiet)
	cFBconfig.ignore_alpha = C.bool(fbConf.IgnoreAlpha)
	cFBconfig.halign = C.uint8_t(fbConf.Halign)
	cFBconfig.valign = C.uint8_t(fbConf.Valign)
	return cFBconfig
}

// Version gets the fbink version
func Version() string {
	vers := C.GoString(C.fbink_version())
	return vers
}

// Open "opens the framebuffer device and returns its fd"
// (from "fbink.h")
func Open() int {
	var resultC C.int
	resultC = C.fbink_open()
	return int(resultC)
}

// Close unmaps the framebuffer and closes the file descripter
func Close(fbfd int) error {
	fdC := C.int(fbfd)
	res := CexitCode(C.fbink_close(fdC))
	return createError(res)
}

// Init initializes the fbink global variables
// See "fbink.h" for detailed usage and explanation
func Init(fbfd int, cfg FBInkConfig) error {
	fbConf := fbconfigGoToC(cfg)
	fdC := C.int(fbfd)
	res := CexitCode(C.fbink_init(fdC, &fbConf))
	return createError(res)
}

// Print prints a string to the screen
// See "fbink.h" for detailed usage and explanation
func Print(fbfd int, str string, cfg FBInkConfig) (int, error) {
	fbConf := fbconfigGoToC(cfg)
	fdC := C.int(fbfd)
	strC := C.CString(str)
	defer C.free(unsafe.Pointer(strC))
	rows := int(C.fbink_print(fdC, strC, &fbConf))
	return rows, createError(CexitCode(rows))
}

// Refresh provides a way of refreshing the eink screen
// See "fbink.h" for detailed usage and explanation
func Refresh(fbfd int, top, left, width, height uint32, waveMode string, blackFlash bool) error {
	fdC := C.int(fbfd)
	topC := C.uint32_t(top)
	leftC := C.uint32_t(left)
	widthC := C.uint32_t(width)
	heightC := C.uint32_t(height)
	waveModeC := C.CString(waveMode)
	defer C.free(unsafe.Pointer(waveModeC))
	blackFlashC := C.bool(blackFlash)
	res := CexitCode(C.fbink_refresh(fdC, topC, leftC, widthC, heightC, waveModeC, blackFlashC))
	return createError(res)
}

// IsFBquirky tests for a quirky framebuffer state
// See "fbink.h" for detailed usage and explanation
func IsFBquirky() bool {
	var resultC C.bool
	resultC = C.fbink_is_fb_quirky()
	return bool(resultC)
}

// PrintProgressBar displays a full width progress bar
// See "fbink.h" for detailed usage and explanation
func PrintProgressBar(fbfd int, percentage uint8, cfg FBInkConfig) error {
	fdC := C.int(fbfd)
	percentC := C.uint8_t(percentage)
	cfgC := fbconfigGoToC(cfg)
	res := CexitCode(C.fbink_print_progress_bar(fdC, percentC, &cfgC))
	return createError(res)
}

// PrintImage will print an image to the screen
// See "fbink.h" for detailed usage and explanation
func PrintImage(fbfd int, imgPath string, targX, targY int16, cfg FBInkConfig) error {
	fdC := C.int(fbfd)
	imgPathC := C.CString(imgPath)
	defer C.free(unsafe.Pointer(imgPathC))
	xC := C.short(targX)
	yC := C.short(targY)
	fbConf := fbconfigGoToC(cfg)
	res := CexitCode(C.fbink_print_image(fdC, imgPathC, xC, yC, &fbConf))
	return createError(res)
}

// ButtonScan will scann for the 'Connect' button on the Kobo USB connect screen
// See "fbink.h" for detailed usage and explanation
func ButtonScan(fbfd int, pressButton, noSleep bool) error {
	fdC := C.int(fbfd)
	pressBtnC := C.bool(pressButton)
	noSleepC := C.bool(noSleep)
	res := CexitCode(C.fbink_button_scan(fdC, pressBtnC, noSleepC))
	return createError(res)
}
