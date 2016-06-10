package main

import "github.com/veandco/go-sdl2/sdl"

type Layer struct {
	Name string
	Geometry
	Rect    sdl.Rect
	Surface *sdl.Surface
	Items   []*Drawable
}

func (L *Layer) AddItem(item Drawable) {
	L.Items = append(L.Items, &item)
}

func (L *Layer) Draw(s *sdl.Surface) {
	for _, item := range L.Items {
		(*item).Draw(L.Surface)
	}
	L.Surface.Blit(&L.Rect, s, &L.Rect)
	L.Surface.FillRect(&L.Rect, 0x00000000)
}
