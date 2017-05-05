package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/iwat/qxargs/internal"

	"github.com/nsf/termbox-go"
)

type _Console struct {
	grepper *internal.Grepper
	query   string
	results []string
}

func newConsole() *_Console {
	return &_Console{
		grepper: internal.NewGrepper(),
	}
}

func (c *_Console) update(query string) []string {
	queryArgs := strings.Split(query, " ")
	findArgs := []string(nil)
	grepArgs := []string(nil)

	for _, arg := range queryArgs {
		if strings.HasPrefix(arg, "?") {
			grepArgs = append(grepArgs, arg[1:])
		} else {
			findArgs = append(findArgs, arg)
		}
	}

	results := []string(nil)
	finder := internal.NewFinder(findArgs...)
	for result := range finder.Channel() {
		matched, _ := c.grepper.Grep(result, grepArgs...)
		if matched {
			results = append(results, result)
		}
		if len(results) >= 10 {
			finder.Reset()
			break
		}
	}

	c.query = query
	c.results = results
	return c.results
}

func (c *_Console) loop(commandArgs []string) ([]string, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	defer termbox.Close()

	current := 0
	command := "qxargs " + strings.Join(commandArgs, " ") + " --"

	for {
		if len(c.results) > 10 {
			c.results = c.results[:10]
		}

		y := 0
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		setCells(0, y, "Select one or all files to be executed:", termbox.ColorDefault, termbox.ColorDefault)
		y += 2

		for i, result := range c.results {
			color := termbox.ColorDefault
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
			{"[TAB/↑/↓]", "select"},
			{"[Enter]", "execute"},
			{"[Ctrl-A]", "execute all"},
			{"[Ctrl-C]", "stop"},
		}

		x := 0
		for _, menu := range menus {
			setCells(x, y, menu.key, termbox.ColorYellow, termbox.ColorDefault)
			x += utf8.RuneCountInString(menu.key) + 1
			setCells(x, y, menu.desc, termbox.ColorDefault, termbox.ColorDefault)
			x += len(menu.desc) + 1
		}
		y++

		setCells(0, y, command, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		setCells(len(command)+1, y, c.query, termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCursor(len(command)+1+len(c.query), y)

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
			case termbox.KeyTab, termbox.KeyArrowDown:
				if len(c.results) > 0 {
					current = (current + 1) % len(c.results)
				}
			case termbox.KeyArrowUp:
				if len(c.results) > 0 {
					current = (current - 1) % len(c.results)
					if current < 0 {
						current = len(c.results) - 1
					}
				}
			case termbox.KeyEnter:
				return []string{c.results[current]}, nil
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if len(c.query) > 0 {
					c.update(c.query[:len(c.query)-1])
					current = 0
				}
			case termbox.KeySpace:
				c.update(c.query + " ")
				current = 0
			default:
				if event.Ch != 0 {
					buf := make([]byte, utf8.UTFMax)
					n := utf8.EncodeRune(buf[:], event.Ch)
					c.update(c.query + string(buf[:n]))
					current = 0
				}
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
