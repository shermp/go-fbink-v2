package main

import (
	"fmt"
	"time"

	"github.com/shermp/go-fbink/gofbink"
)

func main() {
	fbinkOpts := gofbink.FBInkConfig{
		Row:        4,
		Fontmult:   3,
		Fontname:   gofbink.IBM,
		IsCentered: false,
		Valign:     gofbink.Center,
		Halign:     gofbink.Center,
	}
	fb := gofbink.New(&fbinkOpts)
	fb.Open()
	fb.Init()
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("Test line %d", i)
		fb.Println(s)
		time.Sleep(500 * time.Millisecond)
	}
	fb.PrintLastLn("This should update the last line!")
	time.Sleep(1 * time.Second)
	for i := 0; i < 100; i += 10 {
		fb.PrintProgressBar(uint8(i), &fbinkOpts)
		time.Sleep(500 * time.Millisecond)
	}
	fb.Close()
}
