package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli"
)

func probeFile(inFile string) int64 {
	command := fmt.Sprintf("ffprobe -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 %s", inFile)
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

func convert(inFile string, outFile string, args string, totalFrames int64) error {
	var errb bytes.Buffer

	if _, err := os.Stat(inFile); err != nil {
		fmt.Printf("File %s doesn't exists\n", inFile)
		return nil
	}

	outDir := filepath.Dir(outFile)

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist\n", outDir)
		return nil
	}

	if _, err := os.Stat(outFile); err == nil {
		fmt.Printf("File %s already exists\n", outFile)
		return nil
	}

	if len(args) > 0 {
		args = fmt.Sprintf(" %s", args)
	}

	finalArgs := fmt.Sprintf("-i %s%s -progress - -nostats -v error %s", inFile, args, outFile)

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

	//desc := fmt.Sprintf("%-50s", inFile)
	fmt.Println(inFile)

	bar := progressbar.NewOptions64(totalFrames,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		//progressbar.OptionSetDescription(desc),
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
			Name:     "i",
			Usage:    "input file",
			Required: true,
		},
		cli.StringFlag{
			Name:     "o",
			Usage:    "output file",
			Required: true,
		},
		cli.StringFlag{
			Name:  "ffmpeg-args",
			Usage: "arguments to pass onto ffmpeg",
		},
	}

	app.Action = func(c *cli.Context) error {
		inFile := c.String("i")
		outFile := c.String("o")
		args := c.String("ffmpeg-args")

		return convert(inFile, outFile, args, probeFile(inFile))
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
