package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type Drawable interface {
	Draw(*sdl.Surface)
	SetRect(*sdl.Rect)
	GetRect() *sdl.Rect
}

type Rect struct {
	Rect  *sdl.Rect
	Color uint32
}

func (item *Rect) Draw(s *sdl.Surface) {
	log.Println(item.Rect)
	s.FillRect(item.Rect, item.Color)
}

func (item *Rect) GetRect() *sdl.Rect {
	return item.Rect
}

func (item *Rect) SetRect(rect *sdl.Rect) {
	item.Rect = rect
}

type Text struct {
	Rect *sdl.Rect
	Text string
}
