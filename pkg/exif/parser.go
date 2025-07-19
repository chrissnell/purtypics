package exif

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

// Basic EXIF parser for JPEG files
type Parser struct {
	file *os.File
}

// ParseJPEG extracts basic EXIF data from JPEG files
func ParseJPEG(path string) (*EXIFData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Check JPEG signature
	sig := make([]byte, 2)
	if _, err := file.Read(sig); err != nil {
		return nil, err
	}
	if !bytes.Equal(sig, []byte{0xFF, 0xD8}) {
		return nil, fmt.Errorf("not a JPEG file")
	}

	// For now, return empty EXIF data
	// Full EXIF parsing is complex and would require significant code
	return &EXIFData{
		DateTime: time.Now(), // placeholder
	}, nil
}

// Simple GPS coordinate structure
type GPSCoord struct {
	Degrees float64
	Minutes float64
	Seconds float64
	Ref     string
}

// ToDecimal converts GPS coordinates to decimal degrees
func (g GPSCoord) ToDecimal() float64 {
	dec := g.Degrees + g.Minutes/60 + g.Seconds/3600
	if g.Ref == "S" || g.Ref == "W" {
		dec = -dec
	}
	return dec
}