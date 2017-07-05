package timeline

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	Offset = 40960
)

var (
	Version        string
	Location       *time.Location
	TimeEnd        = time.Now()
	TimeStart      = TimeEnd.Add(time.Duration(-1**Minutes) * time.Minute)
	TimeLayout     = flag.String("layout", time.RFC3339, "datetime layout")
	TimeStartValue = flag.String("datetime", "", "start datetime")
	LocationName   = flag.String("location", "Asia/Tokyo", "datetime location name")
	Minutes        = flag.Int64("m", 5, "time minutes")
	isVersion      = flag.Bool("v", false, "show version and exit")
	isHelp         = flag.Bool("h", false, "this help")
)

type TimeLine interface {
	ToConut(line string) error
	ToMetrics() map[string]interface{}
	ToGraph() map[string]mackerelplugin.Graphs
	ParseTime(line string) time.Time
}

type Plugin struct {
	TimeLine
	FileName string
}

func NewPlugin(t TimeLine) Plugin {
	return Plugin{
		TimeLine: t,
	}
}

func (pl Plugin) GraphDefinition() map[string]mackerelplugin.Graphs {
	return pl.ToGraph()
}

func (pl Plugin) FetchMetrics() (map[string]interface{}, error) {
	f, err := os.Open(pl.FileName)

	defer f.Close()

	if err != nil {
		return nil, err
	}

	var r *bufio.Reader
	var i int64

L:
	for {
		i++

		f.Seek((Offset*i)*-1, io.SeekEnd)
		r = bufio.NewReader(f)

		info, err := f.Stat()

		if err != nil {
			return nil, err
		}

		if info.Size() <= Offset*i {
			f.Seek(0, io.SeekStart)
			r = bufio.NewReader(f)

			for {
				b, err := r.ReadBytes('\n')

				if err != nil {
					if err == io.EOF {
						break L
					}

					return nil, err
				}

				if err := pl.ToConut(string(b)); err != nil {
					break L
				}
			}
		}

		// NewLine
		if _, err := r.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				continue
			}

			return nil, err
		}

		if b, err := r.ReadBytes('\n'); err == nil {
			lineTime := pl.ParseTime(string(b))

			// ReSeek
			if TimeStart.Unix() <= lineTime.Unix() {
				continue
			}

			if err := pl.ToConut(string(b)); err != nil {
				break L
			}

			for {
				b, err := r.ReadBytes('\n')

				if err != nil {
					if err == io.EOF {
						break L
					}

					return nil, err
				}

				if err := pl.ToConut(string(b)); err != nil {
					break L
				}
			}
		} else if err != io.EOF {
			return nil, err
		}
	}

	return pl.ToMetrics(), nil
}

func (pl Plugin) Run() error {
	flag.Parse()
	args := flag.Args()

	if *isVersion {
		if len(Version) > 0 {
			fmt.Println("v" + Version)
		}

		return nil
	}

	if *isHelp {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file\n", os.Args[0])
		flag.PrintDefaults()
		return nil
	}

	var err error

	if Location, err = time.LoadLocation(*LocationName); err != nil {
		return err
	}

	if len(*TimeStartValue) > 0 {
		if TimeStart, err = time.ParseInLocation(*TimeLayout, *TimeStartValue, Location); err != nil {
			return err
		}

		TimeEnd = TimeStart.Add(time.Duration(*Minutes) * time.Minute)
	} else {
		TimeEnd = TimeEnd.In(Location)
		TimeStart = TimeEnd.Add(time.Duration(-1**Minutes) * time.Minute)
	}

	if len(args) == 0 {
		return errors.New("file not found")
	}

	pl.FileName = args[0]

	if TimeStart.Unix() > TimeEnd.Unix() {
		return errors.New("End time is small")
	}

	fmt.Fprintf(os.Stderr, "%s -> %s\n", TimeStart.Format(*TimeLayout), TimeEnd.Format(*TimeLayout))

	helper := mackerelplugin.NewMackerelPlugin(pl)
	helper.Run()

	return nil
}
