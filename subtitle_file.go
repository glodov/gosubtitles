package gosubtitles

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

// SubtitleFile is a structure which contains rows of the files,
// so captions and methods
type SubtitleFile struct {
	Rows []Subtitle
}

// Load the path subtitle file into structure
func (me *SubtitleFile) Load(path string) []Subtitle {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 0:00:13.500,0:00:16.200
	// 0:17:37.465,0:17:39.335

	reg, _ := regexp.Compile(`\d+:\d{2}:\d{2}\.\d{3}`)

	current := Subtitle{TimeFrom: 0, TimeTo: 0, Text: ""}

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if reg.MatchString(line) {
			if 0 != current.TimeFrom {
				current.trim()
				me.Rows = append(me.Rows, current)
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
		me.Rows = append(me.Rows, current)
	}

	return me.Rows
}
