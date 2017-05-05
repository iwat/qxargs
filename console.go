package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/iwat/qxargs/internal"

	"github.com/nsf/termbox-go"
)

type _Console struct {
	command string
	engine  *internal.Engine
	query   string

	current int
	results []string
}

func newConsole(commandArgs []string, engine *internal.Engine) *_Console {
	return &_Console{
		command: "qxargs " + strings.Join(commandArgs, " ") + " --",
		engine:  engine,
	}
}

func (c *_Console) update(query string) {
	queryArgs := strings.Split(query, " ")
	c.results = c.engine.Query(queryArgs...)
	c.query = query

	if len(c.results) == 0 {
		c.current = -1
	} else {
		c.current = 0
	}
}

func (c *_Console) moveDown() {
	if len(c.results) > 0 {
		c.current = (c.current + 1) % len(c.results)
	}
}

func (c *_Console) moveUp() {
	if len(c.results) == 0 {
		return
	}
	c.current = (c.current - 1) % len(c.results)
	if c.current < 0 {
		c.current = len(c.results) - 1
	}
}

func (c *_Console) draw() {
	y := 0
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	setCells(0, y, "Select one or all files to be executed:", termbox.ColorDefault, termbox.ColorDefault)
	y += 2

	for i, result := range c.results {
		color := termbox.ColorDefault
		if i == c.current {
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

	setCells(0, y, c.command, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	setCells(len(c.command)+1, y, c.query, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(len(c.command)+1+len(c.query), y)

	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

func (c *_Console) loop() ([]string, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	defer termbox.Close()

	for {
		c.draw()

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
				c.moveDown()
			case termbox.KeyArrowUp:
				c.moveUp()
			case termbox.KeyEnter:
				if len(c.results) > 0 {
					return []string{c.results[c.current]}, nil
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				c.updateBackspace()
			case termbox.KeySpace:
				c.update(c.query + " ")
			default:
				if event.Ch != 0 {
					buf := make([]byte, utf8.UTFMax)
					n := utf8.EncodeRune(buf[:], event.Ch)
					c.update(c.query + string(buf[:n]))
				}
			}
			break
		}
	}
}

func (c *_Console) updateBackspace() {
	if len(c.query) > 0 {
		c.update(c.query[:len(c.query)-1])
	}
}

func setCells(x, y int, s string, fg, bg termbox.Attribute) {
	for _, c := range s {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
