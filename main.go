package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

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
	games := GetGames()
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	// settings := app.Modes[app.Mode].Init()
	cols := 7
	rows := 9
	tw := 230
	th := 107
	w := tw * cols
	h := th * rows
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

	r := 0
	c := 0
	for i, game := range games {
		if c == cols-1 {
			r++
		}
		c = i % cols
		go func(i int, r int, c int, game Game) {
			GetImage(game.Appid)
			if i >= rows*cols {
				return
			}
			layerName := fmt.Sprintf("game_%v", game.Appid)
			image := NewImage(&sdl.Rect{int32(c * tw), int32(r * th), int32(tw), int32(th)}, fmt.Sprintf("cache/%v.bmp", game.Appid), game.Name)
			app.Scene.Lock()
			l, _ := app.Scene.AddLayer(layerName)
			l.AddItem(&image)
			app.Scene.Unlock()
			// fmt.Println(layerName)
		}(i, r, c, game)
	}
	go app.Scene.Run()

	running := true
	m := &sync.Mutex{}
	for running {
		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			fmt.Print(".")
			ret := 1
			switch t := event.(type) {
			case *sdl.QuitEvent:
				ret = 0
			case *sdl.KeyDownEvent:
				// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%s\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				// t.Timestamp, t.Type, sdl.GetScancodeName(t.Keysym.Scancode), t.Keysym.Mod, t.State, t.Repeat)
				key := sdl.GetScancodeName(t.Keysym.Scancode)
				// log.Println(key)
				if t.Keysym.Sym == sdl.K_ESCAPE || t.Keysym.Sym == sdl.K_CAPSLOCK {
					ret = 0
				}
				if key == "Space" {
					go func(m *sync.Mutex) {
						r := *app.Scene.Layers["game_22320"].Items[0]
						z := 1.2
						m.Lock()
						if r.GetScale() == 1 {
							app.Scene.UpLayer("game_22320")
							r.SetScale(z)
							r.Move(-int32(float64(r.GetRect().W)*(z-1))/2, -int32(float64(r.GetRect().H)*(z-1))/2)
						} else {
							r.Move(int32(float64(r.GetRect().W)*(z-1))/2, int32(float64(r.GetRect().H)*(z-1))/2)
							r.SetScale(1)
						}
						m.Unlock()
					}(m)
				}
			}
			if ret == 0 {
				running = false
			}
		}
	}
	return 0
}

func main() {
	if _, err := os.Stat("cache"); err != nil {
		os.Mkdir("cache", 777)
	}
	STEAMID = flag.String("steamID", "76561198049930669", "Steam user ID")
	flag.Parse()
	app := new(Application)
	os.Exit(app.run())
}
