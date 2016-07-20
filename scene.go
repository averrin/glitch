package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Geometry struct {
	Width  int32
	Height int32
}

type Scene struct {
	sync.Mutex
	App         *Application
	Rect        sdl.Rect
	LayersStack []*Layer
	Layers      map[string]*Layer
	Changed     bool
	Geometry
}

var font *ttf.Font
var boldFont *ttf.Font

//NewScene constructor
func NewScene(app *Application, size Geometry) *Scene {
	scene := new(Scene)
	scene.App = app
	scene.Geometry = size
	scene.Rect = sdl.Rect{0, 0, size.Width, size.Height}
	scene.Layers = map[string]*Layer{}

	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := filepath.Join(cwd, "fonts")
	font, _ = ttf.OpenFont(path.Join(dir, "Fantasque Regular.ttf"), 16)
	boldFont, _ = ttf.OpenFont(path.Join(dir, "Fantasque Bold.ttf"), 16)

	scene.AddLayer("root")
	scene.Draw()
	return scene
}

func (S *Scene) Run() {
	for {
		changed := S.Draw()
		if changed {
			S.App.Window.UpdateSurface()
			S.Changed = false
		}
		sdl.Delay(5)
	}
}

func (S *Scene) Reset() {
	S.Layers = map[string]*Layer{}
	S.LayersStack = []*Layer{}
}

func (S *Scene) UpLayer(layerName string) {
	layer := S.Layers[layerName]
	var n int
	var l *Layer
	for n, l = range S.LayersStack {
		if l == layer {
			break
		}
	}
	S.LayersStack = append(S.LayersStack[:n], S.LayersStack[n+1:]...)
	S.LayersStack = append(S.LayersStack, layer)
}

func (S *Scene) Draw() bool {
	S.Lock()
	changed := S.GetChanged()
	if !changed {
		S.Unlock()
		return changed
	}
	fmt.Println("scene changed")
	S.Clear()
	for _, layer := range S.LayersStack {
		layer.Draw(S.App.Surface)
	}
	S.Unlock()
	return changed
}

func (S *Scene) GetChanged() bool {
	changed := S.Changed
	for _, l := range S.LayersStack {
		ch := l.GetChanged()
		if !changed && ch {
			return true
		}
	}
	return changed
}

func (S *Scene) Clear() {
	S.App.Surface.FillRect(&S.Rect, 0xff242424)
}

func (S *Scene) removeLayer(name string) {
	S.Lock()
	_, ok := S.Layers[name]
	if !ok {
		S.Unlock()
		return
	}
	delete(S.Layers, name)
	for i, l := range S.LayersStack {
		if l.Name == name {
			l.Destroy()
			S.LayersStack = append(S.LayersStack[:i], S.LayersStack[i+1:]...)
			break
		}
	}
	S.Changed = true
	S.Unlock()
}

func (S *Scene) AddLayer(name string) (*Layer, error) {
	var layer *Layer
	layer, ok := S.Layers[name]
	if ok {
		return layer, errors.New("Use another layer name")
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
	layer.Surface.FillRect(&layer.Rect, 0x00000000)
	S.LayersStack = append(S.LayersStack, layer)
	S.Changed = true
	return layer, nil
}
