package main

import (
	// Standard library
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	// Image
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	// Video
	"github.com/cbsinteractive/mediainfo"

	// etc
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
)

var (
	searchPath *string

	supportedVideos = []string{
		".avi", ".wmv", ".asf", ".mpg", ".vob", ".mkv", ".dvr-ms", ".mp4", ".mov", ".dat", ".m2ts,", ".mts", ".qt", ".mxf", ".m4v", ".gif",
	}

	supportedImages = []string{
		".jpg", ".bmp", ".png", ".jpeg", ".tiff",
	}
)

func init() {
	// Init search path
	searchPath = flag.String("search", "", "Absolute path to verify assets")
	flag.Parse()
	if *searchPath == "" {
		directory, err := dialog.Directory().Title("Directory to verify assets").Browse()
		if err != nil {
			panic(err)
		}
		*searchPath = directory
	}

	*searchPath = filepath.FromSlash(filepath.Clean(*searchPath))
	if !filepath.IsAbs(*searchPath) {
		logrus.Fatalf("Invalid searchPath(%s)\n", *searchPath)
	}
}

// parseVideo parse video format. TODO
func parseVideo(path string) (*VideoFormat, error) {
	path = filepath.Clean(path)
	i, err := mediainfo.New(path)
	// exit status 3221225477 if mediainfo not installed
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
		Width:     i.VideoTracks[0].Width.Val,
		Height:    i.VideoTracks[0].Height.Val,
		Framerate: i.VideoTracks[0].FrameRate.Val,
		Bitrate:   i.VideoTracks[0].Bitrate.Val,
	}, nil
}

// parseImage parse image format.
func parseImage(path string) (*ImageFormat, error) {
	path = filepath.Clean(path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		im, imgerr := jpeg.Decode(f)
		if imgerr != nil {
			return nil, imgerr
		}
		img = im

	case ".png":
		im, imgerr := png.Decode(f)
		if imgerr != nil {
			return nil, imgerr
		}
		img = im

	case ".tiff":
		im, imgerr := tiff.Decode(f)
		if imgerr != nil {
			return nil, imgerr
		}
		img = im

	case ".bmp":
		im, imgerr := bmp.Decode(f)
		if imgerr != nil {
			return nil, imgerr
		}
		img = im

	default:
		return nil, fmt.Errorf("Unsupported image format")
	}

	return &ImageFormat{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}, nil
}

func init() {
	// logrus.SetLevel(logrus.DebugLevel)

	_, err := exec.LookPath("mediainfo")
	if err != nil {
		logrus.Fatal("Failed to reach mediainfo executable.")
	}
}

func main() {
	// Start
	logrus.Infoln("STARTING...")
	logrus.Infof("searchPath : %s\n", *searchPath)

	// Locate all assets
	wg := &sync.WaitGroup{}
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
					wg.Add(1)
					go func(path string) {
						defer wg.Done()
						logrus.Debugf("Process %s as ext %s file...", entry.Name(), ext)
						f, err := parseVideo(path)
						if err != nil {
							logrus.Warnf("Failed to parse %s:%v\n", path, err)
							return
						}
						if err := f.IsRecommendedHDFormat(); err != nil {
							logrus.Warnf("File %s is not recommended HD format :%v\n", entry.Name(), err)
						} else {
							logrus.Infof("File %s is recommended HD format!\n", entry.Name())
						}
					}(path)
				}
			}

			for _, v := range supportedImages {
				if v == ext {
					wg.Add(1)
					go func(path string) {
						defer wg.Done()
						logrus.Debugf("Process %s as ext %s file...", entry.Name(), ext)
						f, err := parseImage(path)
						if err != nil {
							logrus.Warnf("Failed to parse %s:%v\n", path, err)
							return
						}
						// ???
						if f == nil {
							logrus.Warnf("File %s is unknown image format :%v\n", entry.Name(), err)
							return
						}
						if err := f.IsRecommendedHDFormat(); err != nil {
							logrus.Warnf("File %s is not recommended HD format :%v\n", entry.Name(), err)
						} else {
							logrus.Infof("File %s is recommended HD format!\n", entry.Name())
						}
					}(path)
				}
			}

			return nil
		})
	wg.Wait()
	if err != nil {
		logrus.Fatalln("Failed on WalkDir:", err)
	}
	logrus.Infoln("Press enter to continue")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
