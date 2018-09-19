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

// #cgo LDFLAGS: -L${SRCDIR}/../fbinklib -lfbink
// #include <stdlib.h>
// #include <errno.h>
// #include "../FBInk/fbink.h"
import "C"
import (
	"container/list"
	"errors"
	"fmt"
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
	eNoData     = CexitCode(C.ENODATA) * -1
	eTime       = CexitCode(C.ETIME) * -1
)

// FBFDauto is the automatic fbfd handler
const FBFDauto = int(C.FBFD_AUTO)

// const exitSuccess = int(C.EXIT_SUCCESS)

type FBInkState struct {
	ViewWidth      uint32
	ViewHeight     uint32
	ScreenWidth    uint32
	ScreenHeight   uint32
	ViewHoriOrigin uint8
	ViewVertOrigin uint8
	ViewVertOffset uint8
	BPP            uint32
	FontW          uint16
	FontH          uint16
	FontSizeMult   uint8
	FontName       string
	GlyphWidth     uint8
	GlyphHeight    uint8
	MaxCols        uint16
	MaxRows        uint16
	IsPerfectFit   bool
	UserHZ         int32
	PenFGcolor     uint8
	PenBGcolor     uint8
}

// FBInkConfig is a struct which configures the behavior of fbink
type FBInkConfig struct {
	Row         int16
	Col         int16
	fontmult    uint8
	fontname    Font
	IsInverted  bool
	IsFlashing  bool
	IsCleared   bool
	isCentered  bool
	Hoffset     int16
	Voffset     int16
	IsHalfway   bool
	IsPadded    bool
	FGcolor     FGcolor
	BGcolor     BGcolor
	IsOverlay   bool
	isVerbose   bool
	isQuiet     bool
	IgnoreAlpha bool
	Halign      Align
	Valign      Align
}

// RestrictedConfig is a struct which configures the options that require
// FBInk to be reinitilized
type RestrictedConfig struct {
	Fontmult   uint8
	Fontname   Font
	IsCentered bool
	IsVerbose  bool
	IsQuiet    bool
}

func createError(retValue CexitCode) error {
	switch retValue {
	case exitFailure:
		return errors.New("EXIT_FAILURE")
	case eNoDev:
		return errors.New("ENODEV")
	case eNotSup:
		return errors.New("ENOTSUP")
	case eNoData:
		return errors.New("ENODATA")
	case eTime:
		return errors.New("ETIME")
	default:
		return nil
	}
}

// FBInk contains the active FBInk seesion
type FBInk struct {
	internCfg        FBInkConfig
	fbfd             C.int
	lines            *list.List
	totalRowsWritten int16
}

// New creates an fbInker pointer which clients can
// use to interact with the eink framebuffer
func New(cfg *FBInkConfig, rCfg *RestrictedConfig) *FBInk {
	f := &FBInk{}
	f.fbfd = C.FBFD_AUTO
	f.UpdateRestricted(cfg, rCfg)
	f.internCfg.Row = 4
	f.internCfg.Col = 1
	f.lines = list.New()
	f.lines.PushBack(" ")
	return f
}

func (f *FBInk) newConfigC(cfg *FBInkConfig) C.FBInkConfig {
	var cfgC C.FBInkConfig
	cfgC.row = C.short(cfg.Row)
	cfgC.col = C.short(cfg.Col)
	cfgC.fontmult = C.uint8_t(cfg.fontmult)
	cfgC.fontname = C.uint8_t(cfg.fontname)
	cfgC.is_inverted = C.bool(cfg.IsInverted)
	cfgC.is_flashing = C.bool(cfg.IsFlashing)
	cfgC.is_cleared = C.bool(cfg.IsCleared)
	cfgC.is_centered = C.bool(cfg.isCentered)
	cfgC.hoffset = C.short(cfg.Hoffset)
	cfgC.voffset = C.short(cfg.Voffset)
	cfgC.is_halfway = C.bool(cfg.IsHalfway)
	cfgC.is_padded = C.bool(cfg.IsPadded)
	cfgC.fg_color = C.uint8_t(cfg.FGcolor)
	cfgC.bg_color = C.uint8_t(cfg.BGcolor)
	cfgC.is_overlay = C.bool(cfg.IsOverlay)
	cfgC.is_verbose = C.bool(cfg.isVerbose)
	cfgC.is_quiet = C.bool(cfg.isQuiet)
	cfgC.ignore_alpha = C.bool(cfg.IgnoreAlpha)
	cfgC.halign = C.uint8_t(cfg.Halign)
	cfgC.valign = C.uint8_t(cfg.Valign)
	return cfgC
}

// UpdateRestricted updates cfg with the values in rCfg, which is
// followed by a call to Init()
func (f *FBInk) UpdateRestricted(cfg *FBInkConfig, rCfg *RestrictedConfig) {
	cfg.fontmult = rCfg.Fontmult
	f.internCfg.fontmult = rCfg.Fontmult
	cfg.fontname = rCfg.Fontname
	f.internCfg.fontname = rCfg.Fontname
	cfg.isCentered = rCfg.IsCentered
	f.internCfg.isCentered = rCfg.IsCentered
	cfg.isQuiet = rCfg.IsQuiet
	f.internCfg.isQuiet = rCfg.IsQuiet
	cfg.isVerbose = rCfg.IsVerbose
	f.internCfg.isVerbose = rCfg.IsVerbose
	f.Init(cfg)
}

// Version gets the fbink version
func (f *FBInk) Version() string {
	vers := C.GoString(C.fbink_version())
	return vers
}

// Open the framebuffer device and stores its fd
func (f *FBInk) Open() {
	// Only open if we haven't already obtained a file descriptor
	if f.fbfd == C.FBFD_AUTO {
		f.fbfd = C.fbink_open()
	}
}

// Close unmaps the framebuffer and closes the file descripter
func (f *FBInk) Close() (err error) {
	err = nil
	// Nothing to do unless we obtained a file descriptor!
	if f.fbfd != C.FBFD_AUTO {
		res := CexitCode(C.fbink_close(f.fbfd))
		err = createError(res)
	}
	return err
}

// Init initializes the fbink global variables
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) Init(cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	res := CexitCode(C.fbink_init(f.fbfd, &cfgC))
	return createError(res)
}

// GetState dumps a lot of FBInk internal variables
func (f *FBInk) GetState(cfg *FBInkConfig, state *FBInkState) {
	cfgC := f.newConfigC(cfg)
	stateC := C.FBInkState{}
	C.fbink_get_state(&cfgC, &stateC)
	state.ViewWidth = uint32(stateC.view_width)
	state.ViewHeight = uint32(stateC.view_height)
	state.ScreenWidth = uint32(stateC.screen_width)
	state.ScreenHeight = uint32(stateC.screen_height)
	state.ViewHoriOrigin = uint8(stateC.view_hori_origin)
	state.ViewVertOrigin = uint8(stateC.view_vert_origin)
	state.ViewVertOffset = uint8(stateC.view_vert_offset)
	state.BPP = uint32(stateC.bpp)
	state.FontW = uint16(stateC.font_w)
	state.FontH = uint16(stateC.font_h)
	state.FontSizeMult = uint8(stateC.fontsize_mult)
	state.FontName = C.GoString(stateC.font_name)
	state.GlyphWidth = uint8(stateC.glyph_width)
	state.GlyphHeight = uint8(stateC.glyph_height)
	state.MaxCols = uint16(stateC.max_cols)
	state.MaxRows = uint16(stateC.max_rows)
	state.IsPerfectFit = bool(stateC.is_perfect_fit)
	state.UserHZ = int32(stateC.user_hz)
	state.PenFGcolor = uint8(stateC.pen_fg_color)
	state.PenBGcolor = uint8(stateC.pen_bg_color)
}

// FBprint prints a string to the screen
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) FBprint(str string, cfg *FBInkConfig) (rows int, err error) {
	cfgC := f.newConfigC(cfg)
	strC := C.CString(str)
	defer C.free(unsafe.Pointer(strC))
	rows = int(C.fbink_print(f.fbfd, strC, &cfgC))
	return rows, createError(CexitCode(rows))
}

// Println prints to the screen in the manner of calling fmt.Println()
// Output appears as a set of scrolling lines
func (f *FBInk) Println(a ...interface{}) (n int, err error) {
	str := fmt.Sprint(a...)
	n = len([]byte(str))
	if f.lines.Len() > 5 {
		l := f.lines.Front()
		f.lines.Remove(l)
	}
	f.lines.PushBack(str)
	fbStr := ""
	f.internCfg.Row = 4
	for line := f.lines.Front(); line != nil; line = line.Next() {
		fbStr = line.Value.(string)
		r, _ := f.FBprint(fbStr, &f.internCfg)
		f.internCfg.Row += int16(r)
	}
	if f.internCfg.Row > f.totalRowsWritten {
		f.totalRowsWritten = f.internCfg.Row
	} else if f.internCfg.Row < f.totalRowsWritten {
		row := f.internCfg.Row
		diff := f.totalRowsWritten - row
		for i := row; i < row+diff; i++ {
			f.internCfg.Row = i
			fbStr = " "
			f.FBprint(fbStr, &f.internCfg)
		}
	}
	return n, err
}

// PrintLastLn replaces the last line in the output, without scrolling
func (f *FBInk) PrintLastLn(a ...interface{}) (n int, err error) {
	str := fmt.Sprint(a...)
	n = len([]byte(str))
	l := f.lines.Back()
	f.lines.Remove(l)
	f.lines.PushBack(str)
	fbStr := ""
	f.internCfg.Row = 4
	for line := f.lines.Front(); line != nil; line = line.Next() {
		fbStr = line.Value.(string)
		r, _ := f.FBprint(fbStr, &f.internCfg)
		f.internCfg.Row += int16(r)
	}
	return n, err
}

// Refresh provides a way of refreshing the eink screen
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) Refresh(top, left, width, height uint32, waveMode string, blackFlash bool) error {
	topC := C.uint32_t(top)
	leftC := C.uint32_t(left)
	widthC := C.uint32_t(width)
	heightC := C.uint32_t(height)
	waveModeC := C.CString(waveMode)
	defer C.free(unsafe.Pointer(waveModeC))
	blackFlashC := C.bool(blackFlash)
	res := CexitCode(C.fbink_refresh(f.fbfd, topC, leftC, widthC, heightC, waveModeC, blackFlashC))
	return createError(res)
}

// IsFBquirky tests for a quirky framebuffer state
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) IsFBquirky() bool {
	var resultC C.bool
	resultC = C.fbink_is_fb_quirky()
	return bool(resultC)
}

// PrintProgressBar displays a full width progress bar
// NOTE: percentage should be a number between 0 - 100
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) PrintProgressBar(percentage uint8, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	percentC := C.uint8_t(percentage)
	res := CexitCode(C.fbink_print_progress_bar(f.fbfd, percentC, &cfgC))
	return createError(res)
}

// PrintActivityBar displays a full width activity bar
// NOTE: progress should be a number between 0 - 19.
//       where 0 enables an infinite activity bar!
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) PrintActivityBar(progress uint8, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	progressC := C.uint8_t(progress)
	res := CexitCode(C.fbink_print_activity_bar(f.fbfd, progressC, &cfgC))
	return createError(res)
}

// PrintImage will print an image to the screen
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) PrintImage(imgPath string, targX, targY int16, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	imgPathC := C.CString(imgPath)
	defer C.free(unsafe.Pointer(imgPathC))
	xC := C.short(targX)
	yC := C.short(targY)
	res := CexitCode(C.fbink_print_image(f.fbfd, imgPathC, xC, yC, &cfgC))
	return createError(res)
}

// ButtonScan will scann for the 'Connect' button on the Kobo USB connect screen
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) ButtonScan(pressButton, noSleep bool) error {
	pressBtnC := C.bool(pressButton)
	noSleepC := C.bool(noSleep)
	res := CexitCode(C.fbink_button_scan(f.fbfd, pressBtnC, noSleepC))
	return createError(res)
}

// WaitForUSBMSprocessing waits for the end of a kobo USBMS session
// It also tries to detect a succesful content import
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) WaitForUSBMSprocessing(forceUnplug bool) error {
	forceUnplugC := C.bool(forceUnplug)
	res := CexitCode(C.fbink_wait_for_usbms_processing(f.fbfd, forceUnplugC))
	return createError(res)
}
