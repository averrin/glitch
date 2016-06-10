package main

import "github.com/veandco/go-sdl2/sdl"

type Drawable interface {
	Draw(*sdl.Surface)
	SetRect(*sdl.Rect)
	GetRect() *sdl.Rect
	IsChanged() bool
	GetLastRect() *sdl.Rect
}

type Rect struct {
	Rect     *sdl.Rect
	LastRect *sdl.Rect
	Color    uint32
	Changed  bool
}

func NewRect(rect *sdl.Rect, color uint32) Rect {
	item := new(Rect)
	item.Rect = rect
	item.Color = color
	item.LastRect = item.Rect
	item.Changed = true
	return *item
}

func (item *Rect) Draw(s *sdl.Surface) {
	s.FillRect(item.Rect, item.Color)
	item.Changed = false
	item.LastRect = item.Rect
}

func (item *Rect) GetLastRect() *sdl.Rect {
	return item.LastRect
}

func (item *Rect) GetRect() *sdl.Rect {
	return item.Rect
}

func (item *Rect) SetRect(rect *sdl.Rect) {
	item.LastRect = item.Rect
	item.Rect = rect
	item.Changed = true
}

func (item *Rect) IsChanged() bool {
	return item.Changed
}

type Text struct {
	Rect *sdl.Rect
	Text string
}
