package main

import (
  "bufio"
  "bytes"
  "fmt"
  "os"
  "regexp"
  "strconv"
  "time"

  // "github.com/davecgh/go-spew/spew"
  "github.com/urfave/cli"
)

func secondsFromString(re *regexp.Regexp, str string) (int) {
  h, _ := strconv.Atoi(re.FindStringSubmatch(str)[1])
  hours := ((h * 60) * 60)
  m, _ := strconv.Atoi(re.FindStringSubmatch(str)[2])
  mins := m * 60
  secs, _ := strconv.Atoi(re.FindStringSubmatch(str)[3])

  return (hours + mins + secs)
}

func niceTimeFromSeconds(seconds int) (string) {
  h := seconds / 3600
  m := (seconds % 3600) / 60
  s := seconds % 60

  return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func produceSummary() {
  var line bytes.Buffer
  totalSeconds := 0
  elapsedSeconds := 0
  currentSeconds := 0
  elapsedTime := "00:00:00"
  timeToGo := "00:00:00"
  percentageComplete := "0.00"
  commonRegex := "([0-9]{2}):([0-9]{2}):([0-9]{2})"
  durationRegex, _ := regexp.Compile(fmt.Sprintf("Duration: %s", commonRegex))
  timeRegex, _ := regexp.Compile(fmt.Sprintf("time=%s", commonRegex))

  go func() {
    for {
      if elapsedSeconds < totalSeconds { elapsedSeconds += 1 }

      if totalSeconds > 0 {
        elapsedTime = niceTimeFromSeconds(elapsedSeconds)
        fmt.Printf("\rElapsed %s, ETA %s, Progress %s%%", elapsedTime, timeToGo, percentageComplete)
      } else {
        fmt.Printf("\r-- waiting for input --")
      }
      time.Sleep(1000 * time.Millisecond)
    }
  }()

  scanner := bufio.NewScanner(os.Stdin)
  scanner.Split(bufio.ScanRunes)

  for scanner.Scan() {
    line.WriteString(scanner.Text())

    if totalSeconds == 0 {
      if durationRegex.MatchString(line.String()) {
        totalSeconds = secondsFromString(durationRegex, line.String())
        line.Reset()
      }
    } else {
      if timeRegex.MatchString(line.String()) {
        currentSeconds = secondsFromString(timeRegex, line.String())
        timeToGo = niceTimeFromSeconds(totalSeconds - currentSeconds)
        percentageComplete = fmt.Sprintf("%.2f", (float32(currentSeconds) / float32(totalSeconds) * 100))
        line.Reset()
      }
    }
  }

  fmt.Printf("\rElapsed %s, ETA 00:00:00, Progress 100%%", elapsedTime)
}

func main() {
  app := cli.NewApp()
  app.Name = "ffprogress"
  app.Usage = "Elapsed time, ETA and progress percentage based on your ffmpeg call."
  app.Version = "0.1.0"

  app.Action = func(c *cli.Context) error {
    produceSummary()
    return nil
  }

  app.Run(os.Args)
}
