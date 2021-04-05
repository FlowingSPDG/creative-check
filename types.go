package main

import (
	"fmt"
)

type VideoFormat struct {
	Width     int64
	Height    int64
	Framerate float64
	Bitrate   float64
}

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

type ImageFormat struct {
	Width  int64
	Height int64
}

func (v ImageFormat) IsRecommendedHDFormat() error {
	if v.Width != 1920 {
		return fmt.Errorf("Invalid width size")
	}
	if v.Height != 1080 {
		return fmt.Errorf("Invalid height size")
	}
	return nil
}
