package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *Scene
}

var games []Game
var offset int

const cols = 7
const rows = 9
const tw = 230
const th = 107

func (app *Application) run() int {
	games = GetGames()
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	// settings := app.Modes[app.Mode].Init()
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
			drawItem(app, game, c, r)
		}(i, r, c, game)
	}
	go app.Scene.Run()

	// go func() {
	// 	for {
	// 		offset++
	// 		redraw(app, rows, cols, 1)
	// 		sdl.Delay(1000)
	// 		offset--
	// 		redraw(app, rows, cols, -1)
	// 		sdl.Delay(1000)
	// 	}
	// }()

	running := true
	offset = 0
	// m := &sync.Mutex{}
	for running {
		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// fmt.Print(".")
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
				if key == "Up" {
					if offset != 0 {
						offset--
						redraw(app, rows, cols, -1)
					}
				}
				if key == "Down" {
					if offset < (len(games)/rows - 1) {
						offset++
						redraw(app, rows, cols, 1)
					}
					// go func(m *sync.Mutex) {
					// 	r := *app.Scene.Layers["game_22320"].Items[0]
					// 	z := 1.2
					// 	m.Lock()
					// 	if r.GetScale() == 1 {
					// 		app.Scene.UpLayer("game_22320")
					// 		r.SetScale(z)
					// 		r.Move(-int32(float64(r.GetRect().W)*(z-1))/2, -int32(float64(r.GetRect().H)*(z-1))/2)
					// 	} else {
					// 		r.Move(int32(float64(r.GetRect().W)*(z-1))/2, int32(float64(r.GetRect().H)*(z-1))/2)
					// 		r.SetScale(1)
					// 	}
					// 	m.Unlock()
					// }(m)
				}
			}
			if ret == 0 {
				running = false
			}
		}
	}
	return 0
}

func drawItem(app *Application, game Game, c int, r int) {
	layerName := fmt.Sprintf("game_%v", game.Appid)
	image := NewImage(&sdl.Rect{int32(c * tw), int32(r * th), int32(tw), int32(th)}, fmt.Sprintf("cache/%v.bmp", game.Appid), game.Name)
	app.Scene.Lock()
	l, _ := app.Scene.AddLayer(layerName)
	l.Desc = game.Name
	l.AddItem(&image)
	app.Scene.Unlock()
}

func redraw(app *Application, rows int, cols int, d int) {
	gw := games[offset*cols : (rows+offset+1)*cols]
	gwIDs := func() []string {
		r := []string{}
		for _, g := range gw {
			r = append(r, fmt.Sprintf("game_%v", g.Appid))
		}
		return r
	}()
	counter := 0
	lays := make([]*Layer, len(app.Scene.LayersStack))
	copy(lays, app.Scene.LayersStack)
	for _, l := range lays {
		if l.Name == "root" {
			continue
		}
		b := false
		for _, n := range gwIDs {
			if l.Name == n {
				item := *l.Items[0]
				item.Move(0, int32(-th*d))
				// rect := item.GetRect()
				// fmt.Println("m", l.Desc, rect.X/rect.W, rect.Y/rect.H)
				b = true
				counter++
				break
			}
		}
		if b {
			continue
		}
		// fmt.Println("r", l.Desc)
		app.Scene.removeLayer(l.Name)
	}
	c := 0
	r := rows - 1
	gs := games[(offset+rows)*cols : (offset+rows+1)*cols]
	if d == -1 {
		r = 0
		gs = games[offset*cols : (offset+1)*cols]
	}
	for i, g := range gs {
		c = i % cols
		go func(row int, col int, game Game) {
			// fmt.Println("d", game.Name, c, r)
			drawItem(app, game, col, row)
		}(r, c, g)
	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if _, err := os.Stat("cache"); err != nil {
		os.Mkdir("cache", 0777)
	}
	STEAMID = flag.String("steamID", "76561198049930669", "Steam user ID")
	flag.Parse()
	app := new(Application)
	os.Exit(app.run())
}
