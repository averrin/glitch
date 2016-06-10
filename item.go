package main

import "github.com/veandco/go-sdl2/sdl"

type Drawable interface {
	Draw(*sdl.Surface)
	Move(int32, int32)
	MoveTo(int32, int32)
	IsChanged() bool
	Clear(*sdl.Surface)
	GetRect() *sdl.Rect
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

func (item *Rect) Clear(s *sdl.Surface) {
	s.FillRect(item.LastRect, 0x00000000)
}

func (item *Rect) Move(x int32, y int32) {
	item.LastRect = &sdl.Rect{item.Rect.X, item.Rect.Y, item.Rect.W, item.Rect.H}
	item.Rect.X += x
	item.Rect.Y += y
	item.Changed = true
}

func (item *Rect) MoveTo(x int32, y int32) {
	item.LastRect = &sdl.Rect{item.Rect.X, item.Rect.Y, item.Rect.W, item.Rect.H}
	item.Rect.X = x
	item.Rect.Y = y
	item.Changed = true
}

type Text struct {
	Rect *sdl.Rect
	Text string
}
