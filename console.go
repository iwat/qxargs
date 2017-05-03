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

	c.args = queryArgs
	c.results = results
	return c.results
}

func (c *_Console) loop(commandArgs []string) ([]string, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	defer termbox.Close()

	if len(c.results) > 10 {
		c.results = c.results[:10]
	}

	current := 0
	command := "qxargs " + strings.Join(commandArgs, " ") + " --"
	args := strings.Join(c.args, " ")

	for {
		y := 0
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		setCells(0, y, "Select one or all files to be executed:", termbox.ColorWhite, termbox.ColorBlack)
		y += 2

		for i, result := range c.results {
			color := termbox.ColorWhite
			if i == current {
				color |= termbox.AttrBold
			}
			termbox.SetCell(1, y, '-', color, termbox.ColorDefault)
			setCells(3, y, result, color, termbox.ColorDefault)
			y++
		}
		y += 2

		menus := []struct {
			key  string
			desc string
		}{
			{"[TAB]", "select"},
			{"[Enter]", "execute"},
			{"[Ctrl-A]", "execute all"},
			{"[Ctrl-C]", "stop"},
		}

		x := 0
		for _, menu := range menus {
			setCells(x, y, menu.key, termbox.ColorYellow, termbox.ColorDefault)
			x += len(menu.key) + 1
			setCells(x, y, menu.desc, termbox.ColorWhite, termbox.ColorDefault)
			x += len(menu.desc) + 1
		}
		y++

		setCells(0, y, command, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
		setCells(len(command)+1, y, args, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCursor(len(command)+1+len(args), y)

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
