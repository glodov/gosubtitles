package gosubtitles

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SubtitleFormat is default file subtitle format
var SubtitleFormat = "sbv"

// Subtitle structure
type Subtitle struct {
	TimeFrom float64
	TimeTo   float64
	Text     string
}

// Decode time From and To in Subtitle
func (sub *Subtitle) DecodeTime(from string, to string) {
	sub.TimeFrom = sub.Decode(from)
	sub.TimeTo = sub.Decode(to)
}

func (sub *Subtitle) trim() {
	sub.Text = strings.TrimSpace(sub.Text)
}

// Decode the time from string into float64
func (sub Subtitle) Decode(time string) float64 {
	reg, _ := regexp.Compile(`(\d+):(\d{2}):(\d{2})\.(\d{3})`)
	if !reg.MatchString(time) {
		return -1.0
	}

	matches := reg.FindAllStringSubmatch(time, -1)
	h, _ := strconv.ParseFloat(matches[0][1], 32)
	m, _ := strconv.ParseFloat(matches[0][2], 32)
	s, _ := strconv.ParseFloat(matches[0][3], 32)
	i, _ := strconv.ParseFloat(matches[0][4], 32)

	return h*3600 + m*60 + s + i/1000
}

// Format returns formatted Subtitle as a string
func (sub Subtitle) Format() string {
	// _, divFrom := math.Modf(sub.TimeFrom)
	// _, divTo := math.Modf(sub.TimeTo)
	sFrom := fmt.Sprintf("%.3f", sub.TimeFrom)
	sTo := fmt.Sprintf("%.3f", sub.TimeTo)
	return "" +
		fmt.Sprintf(
			"%d:%02d:%02d.%s",
			int(math.Floor(sub.TimeFrom/3600)),
			int(math.Mod(math.Floor(sub.TimeFrom/60), 60)),
			int(math.Mod(math.Floor(sub.TimeFrom), 60)),
			sFrom[len(sFrom)-3:]) +
		"," +
		fmt.Sprintf(
			"%d:%02d:%02d.%s",
			int(math.Floor(sub.TimeTo/3600)),
			int(math.Mod(math.Floor(sub.TimeTo/60), 60)),
			int(math.Mod(math.Floor(sub.TimeTo), 60)),
			sTo[len(sTo)-3:]) +
		"\n" +
		sub.Text +
		"\n"
}

func loadSubtitles(path string) []Subtitle {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 0:00:13.500,0:00:16.200
	// 0:17:37.465,0:17:39.335

	reg, _ := regexp.Compile(`\d+:\d{2}:\d{2}\.\d{3}`)

	var lines []Subtitle

	current := Subtitle{TimeFrom: 0, TimeTo: 0, Text: ""}

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if reg.MatchString(line) {
			if 0 != current.TimeFrom {
				current.trim()
				lines = append(lines, current)
				current = Subtitle{TimeFrom: 0, TimeTo: 0, Text: ""}
			}
			found := reg.FindAllString(line, 2)
			current.DecodeTime(found[0], found[1])
			// current.TimeFrom = found[0]
			// current.TimeTo = found[1]
		} else {
			current.Text = current.Text + "\n" + line
		}
		// fmt.Println(i, scanner.Text())
	}
	if 0 != current.TimeFrom {
		current.trim()
		lines = append(lines, current)
	}

	return lines
}
