package main

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

type Geometry struct {
	Width  int32
	Height int32
}

type Scene struct {
	App         *Application
	Rect        sdl.Rect
	LayersStack []*Layer
	Layers      map[string]*Layer
	Geometry
}

//NewScene constructor
func NewScene(app *Application, size Geometry) *Scene {
	scene := new(Scene)
	scene.App = app
	scene.Geometry = size
	scene.Rect = sdl.Rect{0, 0, size.Width, size.Height}
	scene.Layers = map[string]*Layer{}
	scene.AddLayer("root")
	scene.AddLayer("test")
	r := Rect{&sdl.Rect{0, 0, 100, 100}, 0xffff0000}
	g := Rect{&sdl.Rect{50, 50, 100, 100}, 0xff00ff00}
	scene.Layers["root"].AddItem(&r)
	scene.Layers["test"].AddItem(&g)
	scene.Draw()
	return scene
}

func (S *Scene) Draw() {
	S.Clear()
	for _, layer := range S.LayersStack {
		layer.Draw(S.App.Surface)
	}
}

func (S *Scene) Clear() {
	S.App.Surface.FillRect(&S.Rect, 0xff242424)
}

func (S *Scene) AddLayer(name string) error {
	var layer *Layer
	layer, ok := S.Layers[name]
	if ok {
		return errors.New("Use another layer name")
	}
	layer = new(Layer)
	S.Layers[name] = layer
	layer.Name = name
	layer.Rect = S.Rect
	amask := uint32(0xff000000)
	rmask := uint32(0x00ff0000)
	gmask := uint32(0x0000ff00)
	bmask := uint32(0x000000ff)
	layer.Surface, _ = sdl.CreateRGBSurface(sdl.SWSURFACE, S.Width, S.Height, 32, rmask, gmask, bmask, amask)
	S.LayersStack = append(S.LayersStack, layer)
	return nil
}
