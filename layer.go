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
		i := (*item)
		if i.IsChanged() {
			i.Clear(L.Surface)
			i.Draw(L.Surface)
		}
	}
	L.Surface.Blit(&L.Rect, s, &L.Rect)
}
