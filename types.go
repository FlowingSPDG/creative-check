package main

import (
	"fmt"
)

// VideoFormat Video file format
type VideoFormat struct {
	Width     int
	Height    int
	Framerate float64
	Bitrate   int
}

// IsRecommendedHDFormat Check if its recommended HD video format or not
func (v VideoFormat) IsRecommendedHDFormat() error {
	if v.Width != 1920 {
		return fmt.Errorf("Invalid width size(%d)", v.Width)
	}
	if v.Height != 1080 {
		return fmt.Errorf("Invalid height size(%d)", v.Height)
	}
	if v.Framerate != 60 && v.Framerate != 59.94 {
		return fmt.Errorf("Invalid frame rate(%v)", v.Framerate)
	}
	return nil
}

// ImageFormat Image file format
type ImageFormat struct {
	Width  int
	Height int
}

// IsRecommendedHDFormat Check if its recommended HD video format or not
func (v ImageFormat) IsRecommendedHDFormat() error {
	if v.Width != 1920 {
		return fmt.Errorf("Invalid width size(%d)", v.Width)
	}
	if v.Height != 1080 {
		return fmt.Errorf("Invalid height size(%d)", v.Height)
	}
	return nil
}
