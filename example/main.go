package main

import (
	"fmt"
	"time"

	"github.com/shermp/go-fbink-v2/gofbink"
)

func testPrints(fb *gofbink.FBInk, cfg *gofbink.FBInkConfig) {
	// Test the console like printing feature
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("Test line %d", i)
		fb.Println(s)
		fmt.Println(s)
		time.Sleep(500 * time.Millisecond)
	}
	// Lets test last line replacement next
	fb.PrintLastLn("This should update the last line!")
	time.Sleep(1 * time.Second)
	// And we finish with a nice progress bar :)
	for i := 0; i <= 100; i += 10 {
		fb.PrintProgressBar(uint8(i), cfg)
		fmt.Println("Progress bar @", i, "%")
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// Set a few initial options
	fbinkOpts := gofbink.FBInkConfig{
		Row:    4,
		Valign: gofbink.Center,
		Halign: gofbink.Center,
	}
	// Set some of the "restricted" options.
	// These are separated out because FBInk requires re-init
	// upon being changed
	rOpts := gofbink.RestrictedConfig{
		Fontmult:   3,
		Fontname:   gofbink.IBM,
		IsCentered: false,
	}
	// fb is a pointer to the object we create, which will contain
	// our FBInk session for the lifetime of our program.

	// NOTE: we create ONE instance of this ONLY.
	fb := gofbink.New(&fbinkOpts, &rOpts)
	// Say hello
	fmt.Println("Using FBInk", fb.Version())
	// Optionlly open a file descriptor. If Open() is not called, FBInk
	// will manage this upon every call to FBInk functions
	fb.Open()
	// Init our FBInk session for the first time. This may be called multiple
	// times as needed, however the current gofbink structure means we should
	// only need to call it once.
	fb.Init(&fbinkOpts)
	// Lets try doing a few things with FBInk...
	testPrints(fb, &fbinkOpts)
	// Now we change some restricted options
	rOpts.IsCentered = true
	rOpts.Fontname = gofbink.UNSCII
	// And update those restricted options. This calls Init() for us :)
	fb.UpdateRestricted(&fbinkOpts, &rOpts)
	time.Sleep(1 * time.Second)
	// We'll do those same FBInk stuff as before, but with the new restricted options
	testPrints(fb, &fbinkOpts)
	// If we've made a call to Open(), we need to clean up...
	fb.Close()
}
