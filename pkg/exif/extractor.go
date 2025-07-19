package exif

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// EXIFData contains photo metadata
type EXIFData struct {
	DateTime     time.Time
	Camera       string
	Lens         string
	ISO          int
	Aperture     float64
	ShutterSpeed string
	FocalLength  int
	GPS          *GPSData
	Orientation  int // EXIF orientation value (1-8)
}

// GPSData contains location information
type GPSData struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

// ExtractMetadata reads EXIF data from an image file
func ExtractMetadata(path string) (*EXIFData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return nil, err // not all images have EXIF
	}

	data := &EXIFData{}

	// Extract datetime
	if dt, err := x.DateTime(); err == nil {
		data.DateTime = dt
	}

	// Extract camera info
	if make, err := x.Get(exif.Make); err == nil {
		if model, err := x.Get(exif.Model); err == nil {
			data.Camera = strings.TrimSpace(make.String() + " " + model.String())
		}
	}

	// Extract lens info
	if lens, err := x.Get(exif.LensModel); err == nil {
		data.Lens = strings.TrimSpace(lens.String())
	}

	// Extract ISO
	if iso, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if val, err := iso.Int(0); err == nil {
			data.ISO = val
		}
	}

	// Extract aperture
	if aperture, err := x.Get(exif.FNumber); err == nil {
		if num, denom, _ := aperture.Rat2(0); denom != 0 {
			data.Aperture = float64(num) / float64(denom)
		}
	}

	// Extract shutter speed
	if shutter, err := x.Get(exif.ExposureTime); err == nil {
		if num, denom, _ := shutter.Rat2(0); denom != 0 {
			if denom > num {
				data.ShutterSpeed = fmt.Sprintf("1/%d", denom/num)
			} else {
				data.ShutterSpeed = fmt.Sprintf("%.1fs", float64(num)/float64(denom))
			}
		}
	}

	// Extract focal length
	if focal, err := x.Get(exif.FocalLength); err == nil {
		if num, denom, _ := focal.Rat2(0); denom != 0 {
			data.FocalLength = int(float64(num) / float64(denom))
		}
	}

	// Extract GPS data
	if lat, lon, err := x.LatLong(); err == nil {
		data.GPS = &GPSData{
			Latitude:  lat,
			Longitude: lon,
		}

		// Try to get altitude
		if alt, err := x.Get(exif.GPSAltitude); err == nil {
			if num, denom, _ := alt.Rat2(0); denom != 0 {
				data.GPS.Altitude = float64(num) / float64(denom)
			}
		}
	}

	// Extract orientation
	if orientation, err := x.Get(exif.Orientation); err == nil {
		if val, err := orientation.Int(0); err == nil {
			data.Orientation = val
		} else {
			data.Orientation = 1 // Default to normal orientation
		}
	} else {
		data.Orientation = 1 // Default to normal orientation
	}

	return data, nil
}

// GetOrientation reads only the EXIF orientation from an image file
func GetOrientation(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 1, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return 1, nil // Default to normal orientation if no EXIF
	}

	// Extract orientation
	if orientation, err := x.Get(exif.Orientation); err == nil {
		if val, err := orientation.Int(0); err == nil {
			return val, nil
		}
	}

	return 1, nil // Default to normal orientation
}

