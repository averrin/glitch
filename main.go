package main

import (
	"flag"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *Scene
}

func (app *Application) run() int {
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	// settings := app.Modes[app.Mode].Init()
	w := 800
	h := 600
	window, err := sdl.CreateWindow("Glitch", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	app.Window = window
	defer app.Window.Destroy()
	renderer, err := sdl.CreateRenderer(app.Window, -1, sdl.RENDERER_ACCELERATED)
	surface, err := app.Window.GetSurface()
	if err != nil {
		panic(err)
	}
	app.Renderer = renderer
	app.Surface = surface
	renderer.Clear()
	app.Scene = NewScene(app, Geometry{int32(w), int32(h)})
	renderer.Present()
	sdl.Delay(5)
	app.Window.UpdateSurface()
	go app.Scene.Run()

	go func() {
		var d int32
		d = 2
		for {
			r := *app.Scene.Layers["root"].Items[0]
			rect := r.GetRect()
			r.Move(0, d)
			sdl.Delay(50)
			if rect.Y == 100 {
				d = -2
			}
			if rect.Y == 0 {
				d = 2
			}
		}
	}()

	go func() {
		var d int32
		d = 2
		for {
			r := *app.Scene.Layers["test"].Items[0]
			rect := r.GetRect()
			r.Move(d, d)
			sdl.Delay(50)
			if rect.Y == 100 {
				d = -2
			}
			if rect.Y == 0 {
				d = 2
			}
		}
	}()

	running := true
	for running {
		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			ret := 1
			switch t := event.(type) {
			case *sdl.QuitEvent:
				ret = 0
			case *sdl.KeyDownEvent:
				// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%s\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				// t.Timestamp, t.Type, sdl.GetScancodeName(t.Keysym.Scancode), t.Keysym.Mod, t.State, t.Repeat)
				key := sdl.GetScancodeName(t.Keysym.Scancode)
				// log.Println(key)
				//TODO: make mode switching more robust
				if t.Keysym.Sym == sdl.K_ESCAPE || t.Keysym.Sym == sdl.K_CAPSLOCK {
					ret = 0
				}
				if key == "Down" {
					r := *app.Scene.Layers["root"].Items[0]
					r.Move(0, 2)
				}
				if key == "Up" {
					r := *app.Scene.Layers["root"].Items[0]
					r.Move(0, -2)
				}
				if key == "Left" {
					r := *app.Scene.Layers["root"].Items[0]
					r.Move(-2, 0)
				}
				if key == "Right" {
					r := *app.Scene.Layers["root"].Items[0]
					r.Move(2, 0)
				}
			default:
				// ret = app.Modes[app.Mode].DispatchEvents(event)
			}
			if ret == 0 {
				running = false
			}
		}
	}
	return 0
}

func main() {
	flag.Parse()
	app := new(Application)
	os.Exit(app.run())
}
