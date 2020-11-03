package ffmpeg

import (
	"fmt"
	"strings"
)

type BatchScreenshotOptions struct {
	OutputPath string
	Quality    int
	TimeStamps []float64
	Width      int
	Verbosity  string
}

func (e *Encoder) BatchScreenshot(probeResult VideoFile, options BatchScreenshotOptions) error {
	if options.Verbosity == "" {
		options.Verbosity = "error"
	}
	if options.Quality == 0 {
		options.Quality = 1
	}

	sFrames := make([]string, len(options.TimeStamps))
	for i, ts := range options.TimeStamps {
		sFrames[i] = fmt.Sprintf("lt(prev_pts*TB\\,%.3f)*gte(pts*TB\\,%.3f)", ts, ts)
	}
	args := []string{
		"-hwaccel", "cuda",
		"-v", options.Verbosity,
		"-y",
		"-i", probeResult.Path,
		"-q:v", fmt.Sprintf("%v", options.Quality),
		"-vf", fmt.Sprintf("select='%s',scale=%d:-1", strings.Join(sFrames, "+"), options.Width),
		"-vsync", "drop",
		"-f", "image2",
		options.OutputPath,
	}
	_, err := e.run(probeResult, args)

	return err
}
