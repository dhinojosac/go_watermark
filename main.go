package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdf "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func StampTime() string {
	loc, err := time.LoadLocation("America/Santiago")
	t := time.Now()
	if err == nil {
		t = t.In(loc)
	}
	//fmt.Println(t.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
	return t.Format("02-Jan-2006 | 15:04:05")
}

func hasWatermarks(inFile string) bool {
	ok, err := api.HasWatermarksFile(inFile, nil)
	if err != nil {
		fmt.Printf("Checking for watermarks: %s: %v\n", inFile, err)
	}
	return ok
}

func main() {
	/*
		inputPtr := flag.String("i", "Input pdf file to add watermark", "a string")
		watermarkPtr := flag.String("w", "Text to add as watermark", "a string")
		flag.Parse()
		if *inputPtr == "Input pdf file to add watermark" {
			fmt.Println("[!] Input pdf file not set")
			fmt.Println("watermark.exe  Input pdf file not set")
		}
		if *watermarkPtr == "Text to add as watermark" {
			fmt.Println("[!] Text to add as watermark not set")
			fmt.Println("[!] Input pdf file not set")
		}
		os.Exit(1)
	*/

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter pdf filename: ")
	fn, _ := reader.ReadString('\n')
	fmt.Println(fn)
	fmt.Print("Enter text to add as watermark: ")
	wt, _ := reader.ReadString('\n')
	fmt.Println(wt)

	msg := "StampWaterMark"
	inFile := strings.TrimSpace(strings.Replace(fn, "\n", "", -1))
	watermarkText := wt
	outFile := strings.Replace(inFile, ".pdf", "", -1) + "_wm.pdf"
	onTop := true // we are testing stamps

	// Check for existing stamps.
	if ok := hasWatermarks(inFile); ok {
		fmt.Printf("Watermarks found: %s\n", inFile)
	}

	// Stamp all pages.
	wm, err := pdf.ParseTextWatermarkDetails(watermarkText, "op:0.2, sc:1.0 rel, off: -30 -0", onTop)
	if err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}
	if err := api.AddWatermarksFile(inFile, outFile, nil, wm, nil); err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}

	timestamp := StampTime()
	// Stamp all pages.
	wm, err = pdf.ParseTextWatermarkDetails(timestamp, "op:0.2, sc:0.7 rel, off: 60 -30", onTop)
	if err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}
	if err = api.AddWatermarksFile(outFile, outFile, nil, wm, nil); err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}

	// Check for existing stamps.
	if ok := hasWatermarks(outFile); !ok {
		fmt.Printf("No watermarks found: %s\n", outFile)
	}
	_, err = reader.ReadString('\n')
	fmt.Println(err)

}
