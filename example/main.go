package main

import (
	"fmt"
	"time"

	"github.com/shermp/go-fbink/gofbink"
)

func testPrints(fb *gofbink.FBInk, cfg *gofbink.FBInkConfig) {
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("Test line %d", i)
		fb.Println(s)
		time.Sleep(500 * time.Millisecond)
	}
	fb.PrintLastLn("This should update the last line!")
	time.Sleep(1 * time.Second)
	for i := 0; i < 100; i += 10 {
		fb.PrintProgressBar(uint8(i), cfg)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	fbinkOpts := gofbink.FBInkConfig{
		Row:    4,
		Valign: gofbink.Center,
		Halign: gofbink.Center,
	}

	rOpts := gofbink.RestrictedConfig{
		Fontmult:   3,
		Fontname:   gofbink.IBM,
		IsCentered: false,
	}
	fb := gofbink.New(&fbinkOpts, &rOpts)
	fb.Open()
	fb.Init(&fbinkOpts)
	testPrints(fb, &fbinkOpts)
	rOpts.IsCentered = true
	rOpts.Fontname = gofbink.UNSCII
	fb.UpdateRestricted(&fbinkOpts, &rOpts)
	time.Sleep(1 * time.Second)
	testPrints(fb, &fbinkOpts)
	fb.Close()
}
