/*
	FBInk: FrameBuffer eInker, a tool to print text & images on eInk devices (Kobo/Kindle)
	Copyright (C) 2018-2019 NiLuJe <ninuje@gmail.com>

	go-fbink: A Go wrapper for FBInk
	Copyright (C) 2018-2019 Sherman Perry

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

// #cgo LDFLAGS: -L${SRCDIR}/../fbinklib -lfbink -lm
// #include <stdlib.h>
// #include <errno.h>
// #include "fbink.h"
import "C"
import (
	"container/list"
	"errors"
	"fmt"
	"image"
	"strings"
	"unicode/utf8"
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
	Terminus
	TerminusB
	Fatty
	Spleen
	Tewi
	TewiB
	Topaz
	MicroKnight
	VGA
)

// FontStyle type
type FontStyle int

// FontStyle constants
const (
	FntRegular FontStyle = iota
	FntItalic
	FntBold
	FntBoldItalic
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

// WaveFormMode type
type WaveFormMode uint8

// WaveFormMode constants
const (
	WfmAUTO WaveFormMode = iota
	WfmDU
	WfmGC16
	WfmGC4
	WfmA2
	WfmGL16
	WfmREAGL
	WfmREAGLD
	WfmGC16_FAST
	WfmGL16_FAST
	WfmDU4
	WfmGL4
	WfmGL16_INV
	WfmGCK16
	WfmGLKW16
	WfmINIT
)

// HWDither type
type HWDither uint8

// HWDither constants
const (
	DitherPassthrough HWDither = iota
	DitherFloydSteingberg
	DitherAtkinson
	DitherOrdered
	DitherQuantOnly
)

// NTXRota type
type NTXRota uint8

// NTXRota constants
const (
	NTXRotaStraight NTXRota = iota
	NTXRotaAllInverted
	NTXRotaOddInverted
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
	eInval      = CexitCode(C.EINVAL) * -1
	eIlSeq      = CexitCode(C.EILSEQ) * -1
	eNoSpc      = CexitCode(C.ENOSPC) * -1
)

// FBFDauto is the automatic fbfd handler
const FBFDauto = int(C.FBFD_AUTO)

// const exitSuccess = int(C.EXIT_SUCCESS)

// FBInkState stores a snapshot of some of FBInk's internal variables
type FBInkState struct {
	UserHZ         int
	FontName       string
	ViewWidth      uint32
	ViewHeight     uint32
	ScreenWidth    uint32
	ScreenHeight   uint32
	BPP            uint32
	DeviceName     string
	DeviceCodename string
	DevicePlatform string
	DeviceId       uint16
	PenFGcolor     uint8
	PenBGcolor     uint8
	ScreenDPI      uint16
	FontW          uint16
	FontH          uint16
	MaxCols        uint16
	MaxRows        uint16
	ViewHoriOrigin uint8
	ViewVertOrigin uint8
	ViewVertOffset uint8
	FontSizeMult   uint8
	GlyphWidth     uint8
	GlyphHeight    uint8
	IsPerfectFit   bool
	IsKoboNonMT    bool
	NTXBootRota    uint8
	NTXRotaQuirk   NTXRota
	CurrentRota    uint8
	CanRotate      bool
}

// FBInkConfig is a struct which configures the behavior of fbink
type FBInkConfig struct {
	Row          int16
	Col          int16
	fontmult     uint8
	fontname     Font
	IsInverted   bool
	IsFlashing   bool
	IsCleared    bool
	isCentered   bool
	Hoffset      int16
	Voffset      int16
	IsHalfway    bool
	IsPadded     bool
	IsRpadded    bool
	fgColor      FGcolor
	bgColor      BGcolor
	IsOverlay    bool
	IsBGless     bool
	isFGless     bool
	noViewport   bool
	isVerbose    bool
	isQuiet      bool
	IgnoreAlpha  bool
	Halign       Align
	Valign       Align
	ScaledWidth  int16
	ScaledHeight int16
	WfmMode      WaveFormMode
	IsDithered   bool
	SWDithering  bool
	IsNightmode  bool
	NoRefresh    bool
}

// FBInkOTConfig is a struct which configures OpenType specific options
type FBInkOTConfig struct {
	Margins struct {
		Top    int16
		Bottom int16
		Left   int16
		Right  int16
	}
	SizePt       float32
	SizePx       uint16
	IsCentred    bool
	IsFormatted  bool
	ComputeOnly  bool
	NoTruncation bool
}

type FBInkOTFit struct {
	ComputedLines  uint16
	RenderedLines  uint16
	Truncated      bool
}

type FBInkDump struct {
	data   *uint8
	Size   uint
	x      uint16
	y      uint16
	w      uint16
	h      uint16
	Rota   uint8
	BPP    uint8
	IsFull bool
}

type FBInkRect struct {
	Top    uint16
	Left   uint16
	Width  uint16
	Height uint16
}

// RestrictedConfig is a struct which configures the options that require
// FBInk to be reinitilized
type RestrictedConfig struct {
	Fontmult   uint8
	Fontname   Font
	IsCentered bool
	FGcolor    FGcolor
	BGcolor    BGcolor
	NoViewport bool
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
	case eInval:
		return errors.New("EINVAL")
	case eIlSeq:
		return errors.New("EILSEQ")
	case eNoSpc:
		return errors.New("ENOSPC")
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
	f.internCfg.Row = 1
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
	cfgC.is_rpadded = C.bool(cfg.IsRpadded)
	cfgC.fg_color = C.uint8_t(cfg.fgColor)
	cfgC.bg_color = C.uint8_t(cfg.bgColor)
	cfgC.is_overlay = C.bool(cfg.IsOverlay)
	cfgC.is_bgless = C.bool(cfg.IsBGless)
	cfgC.is_fgless = C.bool(cfg.isFGless)
	cfgC.no_viewport = C.bool(cfg.noViewport)
	cfgC.is_verbose = C.bool(cfg.isVerbose)
	cfgC.is_quiet = C.bool(cfg.isQuiet)
	cfgC.ignore_alpha = C.bool(cfg.IgnoreAlpha)
	cfgC.halign = C.uint8_t(cfg.Halign)
	cfgC.valign = C.uint8_t(cfg.Valign)
	cfgC.scaled_width = C.short(cfg.ScaledWidth)
	cfgC.scaled_height = C.short(cfg.ScaledHeight)
	cfgC.wfm_mode = C.uint8_t(cfg.WfmMode)
	cfgC.is_dithered = C.bool(cfg.IsDithered)
	cfgC.sw_dithering = C.bool(cfg.SWDithering)
	cfgC.is_nightmode = C.bool(cfg.IsNightmode)
	cfgC.no_refresh = C.bool(cfg.NoRefresh)
	return cfgC
}

func (f *FBInk) newOTConfig(otCfg *FBInkOTConfig) C.FBInkOTConfig {
	var otCfgC C.FBInkOTConfig
	otCfgC.margins.top = C.short(otCfg.Margins.Top)
	otCfgC.margins.bottom = C.short(otCfg.Margins.Bottom)
	otCfgC.margins.left = C.short(otCfg.Margins.Left)
	otCfgC.margins.right = C.short(otCfg.Margins.Right)
	otCfgC.size_pt = C.float(otCfg.SizePt)
	otCfgC.size_px = C.uint16_t(otCfg.SizePx)
	otCfgC.is_centered = C.bool(otCfg.IsCentred)
	otCfgC.is_formatted = C.bool(otCfg.IsFormatted)
	otCfgC.compute_only = C.bool(otCfg.ComputeOnly)
	otCfgC.no_truncation = C.bool(otCfg.NoTruncation)
	return otCfgC
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
	cfg.fgColor = rCfg.FGcolor
	f.internCfg.fgColor = rCfg.FGcolor
	cfg.bgColor = rCfg.BGcolor
	f.internCfg.bgColor = rCfg.BGcolor
	cfg.noViewport = rCfg.NoViewport
	f.internCfg.noViewport = rCfg.NoViewport
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

// AddOTfont registers an OpenType or TrueType font with FBInk
// At least one font needs to be specified to use the OT print function
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) AddOTfont(filename string, fntStyle FontStyle) error {
	fnC := C.CString(filename)
	defer C.free(unsafe.Pointer(fnC))
	res := CexitCode(C.fbink_add_ot_font(fnC, C.FONT_STYLE_T(fntStyle)))
	return createError(res)
}

// FreeOTfonts frees any loaded OT font. This MUST be called at the
// conclusion of OT printing, to avoid memory leaks
func (f *FBInk) FreeOTfonts() error {
	res := CexitCode(C.fbink_free_ot_fonts())
	return createError(res)
}

// GetState dumps a lot of FBInk internal variables
func (f *FBInk) GetState(cfg *FBInkConfig, state *FBInkState) {
	cfgC := f.newConfigC(cfg)
	stateC := C.FBInkState{}
	C.fbink_get_state(&cfgC, &stateC)
	state.UserHZ = int(stateC.user_hz)
	state.FontName = C.GoString(stateC.font_name)
	state.ViewWidth = uint32(stateC.view_width)
	state.ViewHeight = uint32(stateC.view_height)
	state.ScreenWidth = uint32(stateC.screen_width)
	state.ScreenHeight = uint32(stateC.screen_height)
	state.BPP = uint32(stateC.bpp)
	state.DeviceName = C.GoString(&stateC.device_name[0])
	state.DeviceCodename = C.GoString(&stateC.device_codename[0])
	state.DevicePlatform = C.GoString(&stateC.device_platform[0])
	state.DeviceId = uint16(stateC.device_id)
	state.PenFGcolor = uint8(stateC.pen_fg_color)
	state.PenBGcolor = uint8(stateC.pen_bg_color)
	state.ScreenDPI = uint16(stateC.screen_dpi)
	state.FontW = uint16(stateC.font_w)
	state.FontH = uint16(stateC.font_h)
	state.MaxCols = uint16(stateC.max_cols)
	state.MaxRows = uint16(stateC.max_rows)
	state.ViewHoriOrigin = uint8(stateC.view_hori_origin)
	state.ViewVertOrigin = uint8(stateC.view_vert_origin)
	state.ViewVertOffset = uint8(stateC.view_vert_offset)
	state.FontSizeMult = uint8(stateC.fontsize_mult)
	state.GlyphWidth = uint8(stateC.glyph_width)
	state.GlyphHeight = uint8(stateC.glyph_height)
	state.IsPerfectFit = bool(stateC.is_perfect_fit)
	state.IsKoboNonMT = bool(stateC.is_kobo_non_mt)
	state.NTXBootRota = uint8(stateC.ntx_boot_rota)
	state.NTXRotaQuirk = NTXRota(stateC.ntx_rota_quirk)
	state.CurrentRota = uint8(stateC.current_rota)
	state.CanRotate = bool(stateC.can_rotate)
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

// PrintOT prints a string to the framebuffer using OpenType or TrueType fonts
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) PrintOT(str string, otCfg *FBInkOTConfig, fbCfg *FBInkConfig) (int, error) {
	fbCfgC := f.newConfigC(fbCfg)
	otCfgC := f.newOTConfig(otCfg)
	strC := C.CString(str)
	defer C.free(unsafe.Pointer(strC))
	res := C.fbink_print_ot(f.fbfd, strC, &otCfgC, &fbCfgC, nil)
	return int(res), createError(CexitCode(res))
}

// Println prints to the screen in the manner of calling fmt.Println()
// Output appears as a set of scrolling lines
func (f *FBInk) Println(a ...interface{}) (n int, err error) {
	str := fmt.Sprint(a...)
	n = len([]byte(str))
	if f.lines.Len() > 8 {
		l := f.lines.Front()
		f.lines.Remove(l)
	}
	f.lines.PushBack(str)
	fbStr := ""
	f.internCfg.Row = 4
	state := FBInkState{}
	f.GetState(&f.internCfg, &state)
	for line := f.lines.Front(); line != nil; line = line.Next() {
		fbStr = line.Value.(string)
		strLen := utf8.RuneCountInString(fbStr)
		numRows := strLen / (int(state.MaxCols) - 1)
		for i := 0; i <= numRows; i++ {
			space := strings.Repeat(" ", (int(state.MaxCols) - 1))
			f.FBprint(space, &f.internCfg)
		}
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
func (f *FBInk) Refresh(top, left, width, height uint32, ditherMode HWDither, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	topC := C.uint32_t(top)
	leftC := C.uint32_t(left)
	widthC := C.uint32_t(width)
	heightC := C.uint32_t(height)
	ditherModeC := C.uint8_t(ditherMode)
	res := CexitCode(C.fbink_refresh(f.fbfd, topC, leftC, widthC, heightC, ditherModeC, &cfgC))
	return createError(res)
}

// // IsFBquirky tests for a quirky framebuffer state
// // See "fbink.h" for detailed usage and explanation
// func (f *FBInk) IsFBquirky() bool {
// 	var resultC C.bool
// 	resultC = C.fbink_is_fb_quirky()
// 	return bool(resultC)
// }

// ReInit handles cases where the framebuffer state such as bit depth
// or rotation may change
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) ReInit(cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	res := CexitCode(C.fbink_reinit(f.fbfd, &cfgC))
	return createError(res)
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

// PrintRawData prints raw scanlines to the screen, without having to save image
// to disk beforehand. Useful for images created programatically.
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) PrintRawData(data []byte, w, h int, xOff, yOff uint16, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	res := CexitCode(C.fbink_print_raw_data(
		f.fbfd,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(w),
		C.int(h),
		C.size_t(len(data)),
		C.short(xOff),
		C.short(yOff),
		&cfgC))
	return createError(res)
}

// PrintRBGA prints an image stored in an image.RGBA
func (f *FBInk) PrintRBGA(xOff, yOff int16, im *image.RGBA, cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	w := im.Rect.Max.X - im.Rect.Min.X
	h := im.Rect.Max.Y - im.Rect.Min.Y
	res := CexitCode(C.fbink_print_raw_data(
		f.fbfd,
		(*C.uchar)(unsafe.Pointer(&im.Pix[0])),
		C.int(w),
		C.int(h),
		C.size_t(im.Stride*h),
		C.short(xOff),
		C.short(yOff),
		&cfgC))
	return createError(res)
}

// ClearScreen simply clears the screen to white
// See "fbink.h" for detailed usage and explanation
func (f *FBInk) ClearScreen(cfg *FBInkConfig) error {
	cfgC := f.newConfigC(cfg)
	res := CexitCode(C.fbink_cls(f.fbfd, &cfgC))
	return createError(res)
}

// TODO: fbink_dump, fbink_region_dump, fbink_restore, fbink_free_dump_data
// TODO: fbink_get_last_rect

// ButtonScan will scan for the 'Connect' button on the Kobo USB connect screen
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
