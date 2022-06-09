package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli"
)

func probeFile(file string) int64 {
	command := fmt.Sprintf("ffprobe -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 %s", file)
	args := strings.Fields(command)

	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		panic(err)
	}

	val, err := strconv.ParseInt(strings.TrimSuffix(string(output), "\n"), 10, 64)
	if err != nil {
		panic(err)
	}

	return val
}

func convert(in_file string, args string, totalFrames int66) error {
	var errb bytes.Buffer

	out_file := fmt.Sprintf("output/%s", in_file)

	if _, err := os.Stat(out_file); errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("%s already exists", out_file)
		panic(msg)
	}

	if len(args) > 0 {
		args = fmt.Sprintf(" %s", args)
	}

	finalArgs := fmt.Sprintf("-i %s%s -progress - -nostats -v error %s", in_file, args, out_file)

	cmd := exec.Command("ffmpeg", strings.Split(finalArgs, " ")...)

	cmd.Stderr = &errb
	output, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	cmd.Start()

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanWords)

	frameRe := regexp.MustCompile(`frame=(.+)`)

	var currentFrame int64 = 0
	var previousFrame int64 = 0

	bar := progressbar.NewOptions64(totalFrames,
		progressbar.OptionSetDescription(in_file),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\n")
		}),
		progressbar.OptionFullWidth(),
	)

	bar.RenderBlank()

	for scanner.Scan() {
		line := scanner.Text()

		match := frameRe.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		currentFrame, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			panic(err)
		}

		if currentFrame > totalFrames {
			bar.Finish()
		} else {
			bar.Add(int(currentFrame - previousFrame))
		}

		previousFrame = currentFrame
	}

	if len(errb.String()) > 0 {
		fmt.Println("err:", errb.String())
	}

	return cmd.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "ffprogress"
	app.Usage = "Elapsed time, ETA and progress percentage based on your ffmpeg call."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "file",
			Usage:    "file to process",
			Required: true,
		},
		cli.StringFlag{
			Name:  "ffmpeg-args",
			Usage: "arguments to pass onto ffmpeg",
		},
	}

	app.Action = func(c *cli.Context) error {
		file := c.String("file")
		args := c.String("ffmpeg-args")

		return convert(file, args, probeFile(file))
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
