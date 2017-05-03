package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

type _Console struct {
	finder  *_Finder
	grepper *_Grepper
	args    []string
	results []string
}

func newConsole(finder *_Finder, grepper *_Grepper) *_Console {
	return &_Console{
		finder:  finder,
		grepper: grepper,
	}
}

func (c *_Console) update(queryArgs []string) []string {
	findArgs := []string(nil)
	grepArgs := []string(nil)

	for _, arg := range queryArgs {
		if strings.HasPrefix(arg, "?") {
			grepArgs = append(grepArgs, arg[1:])
		} else {
			findArgs = append(findArgs, arg)
		}
	}

	results, err := c.finder.Find(findArgs...)
	if err != nil {
		panic(err)
	}

	results, err = c.grepper.Grep(results, grepArgs...)
	if err != nil {
		panic(err)
	}

	c.results = results
	return c.results
}

func (c *_Console) loop() ([]string, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	defer termbox.Close()

	if len(c.results) > 10 {
		c.results = c.results[:10]
	}
	current := 0

	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		setCells(0, 0, "Select one or all files to be executed:", termbox.ColorWhite, termbox.ColorBlack)
		for i, result := range c.results {
			color := termbox.ColorWhite
			if i == current {
				color |= termbox.AttrBold
			}
			termbox.SetCell(1, i+2, '-', color, termbox.ColorDefault)
			setCells(3, i+2, result, color, termbox.ColorDefault)
		}
		setCells(0, len(c.results)+3, "[TAB] to select", termbox.ColorDefault, termbox.ColorDefault)
		setCells(0, len(c.results)+4, "[Enter] to execute", termbox.ColorDefault, termbox.ColorDefault)
		setCells(0, len(c.results)+5, "[Ctrl-A] to execute all, [Ctrl-C] to stop", termbox.ColorDefault, termbox.ColorDefault)

		termbox.Flush()

		for {
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				fmt.Println(event.Type)
				continue
			}

			switch event.Key {
			case termbox.KeyCtrlC:
				return nil, nil
			case termbox.KeyCtrlA:
				return c.results, nil
			case termbox.KeyTab:
				current = (current + 1) % len(c.results)
			case termbox.KeyEnter:
				return []string{c.results[current]}, nil
			}
			break
		}
	}

	return c.results, nil
}

func setCells(x, y int, s string, fg, bg termbox.Attribute) {
	for _, c := range s {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
