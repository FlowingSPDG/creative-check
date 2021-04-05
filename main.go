package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cbsinteractive/mediainfo"
	"github.com/sirupsen/logrus"
)

var (
	searchPath *string

	supportedVideos = []string{
		".avi", ".wmv", ".asf", ".mpg", ".vob", ".mkv", ".dvr-ms", ".mp4", ".mov", ".dat", ".m2ts,", ".mts", ".qt", ".mxf", ".m4v",
	}

	supportedImages = []string{
		"jpg", "bmp", "png", "gif", "jpeg",
	}
)

func init() {
	// Init search path
	searchPath = flag.String("search", "", "Absolute path to search assets")
	flag.Parse()
	*searchPath = filepath.FromSlash(filepath.Clean(*searchPath))
	if !filepath.IsAbs(*searchPath) {
		logrus.Fatalf("Invalid searchPath(%s)\n", *searchPath)
	}
}

// parseVideo parse video format. TODO
func parseVideo(path string) (*VideoFormat, error) {

	path = filepath.Clean(path)
	i, err := mediainfo.New(path)
	// exit status 3221225477?
	if err != nil {
		return nil, err
	}

	if i.VideoTracks == nil {
		return nil, fmt.Errorf("Unknown video type")
	}
	if len(i.VideoTracks) != 1 {
		return nil, fmt.Errorf("Unknown video type")
	}

	return &VideoFormat{
		Width:     int64(i.VideoTracks[0].Width.Val),
		Height:    int64(i.VideoTracks[0].Height.Val),
		Framerate: i.General.FrameRate.Val,
		Bitrate:   float64(i.VideoTracks[0].Bitrate.Val),
	}, nil
}

// parseImage parse image format. TODO
func parseImage(ext, path string) (*VideoFormat, error) {
	return nil, nil
}

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// Start
	logrus.Infoln("STARTING...")
	logrus.Infof("searchPath : %s\n", *searchPath)

	// Locate all assets
	err := filepath.WalkDir(*searchPath,
		func(path string, entry fs.DirEntry, err error) error {
			if entry.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == "" {
				return nil
			}
			for _, v := range supportedVideos {
				if v == ext {
					logrus.Debugf("Process %s as ext %s file...", entry.Name(), ext)
					f, err := parseVideo(path)
					if err != nil {
						logrus.Warnf("Failed to parse %s:%v\n", path, err)
						continue
					}
					if err := f.IsRecommendedHDFormat(); err != nil {
						logrus.Warnf("File %s is not recommended HD format :%v\n", entry.Name(), err)
					} else {
						logrus.Infof("File %s is recommended HD format!\n", entry.Name())
					}
				}
			}

			for _, v := range supportedImages {
				if v == ext {
					logrus.Debugf("Process %s as ext %s file...", path, ext)
				}
			}

			return nil
		})

	if err != nil {
		logrus.Fatalln("Failed on WalkDir:", err)
	}
}
