package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"net/http"
	_ "net/http/pprof"

	"github.com/averrin/seker"
	st "github.com/averrin/shodan/modules/steam"
	"github.com/spf13/viper"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *seker.Scene
}

func (app *Application) GetSurface() *sdl.Surface {
	return app.Surface
}

func (app *Application) GetWindow() *sdl.Window {
	return app.Window
}

type Mover struct {
	Item seker.Drawable
	Y    int32
}

var games []st.Game
var offset int
var LOCK sync.Mutex
var steam st.Steam

const cols = 7
const rows = 9
const tw = 230
const th = 107

func Shuffle(a []Mover) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func (app *Application) run() int {
	LOCK = sync.Mutex{}
	games = steam.GetGames()
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
	app.Scene = seker.NewScene(app, seker.Geometry{int32(w), int32(h)})
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
		go func(i int, r int, c int, game st.Game) {
			GetImage(game.Appid)
			if i >= (rows+1)*cols {
				return
			}
			drawItem(app, game, c, r)
		}(i, r, c, game)
	}
	go app.Scene.Run()

	offset = 0
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	d := 1
	// 	for {
	// 		log.Print(offset, len(games)/cols)
	// 		if offset == len(games)/cols {
	// 			d = -1
	// 		}
	// 		if offset == 0 {
	// 			d = 1
	// 		}
	// 		redraw(app, rows, cols, d)
	// 		time.Sleep(500 * time.Millisecond)
	// 		// redraw(app, rows, cols, -1)
	// 		// time.Sleep(2 * time.Second)
	// 	}
	// }()

	running := true
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
						redraw(app, rows, cols, -1)
					}
				}
				if key == "Down" {
					if offset < (len(games)/rows - 2) {
						redraw(app, rows, cols, 1)
					}
				}
			}
			if ret == 0 {
				running = false
			}
		}
	}
	return 0
}

func drawItem(app *Application, game st.Game, c int, r int) {
	layerName := fmt.Sprintf("game_%v", game.Appid)
	image := seker.NewImage(&sdl.Rect{int32(c * tw), int32(r * th), int32(tw), int32(th)}, fmt.Sprintf("cache/%v.bmp", game.Appid), game.Name)
	l, _ := app.Scene.AddLayer(layerName)
	l.Desc = game.Name
	l.AddItem(&image)
}

func redraw(app *Application, rows int, cols int, d int) {
	LOCK.Lock()
	offset += d
	log.Print(offset)

	WG := sync.WaitGroup{}
	end := (rows + offset + 2) * cols
	if end > len(games)-1 {
		end = len(games) - 1
	}
	start := offset * cols
	if offset > 0 {
		start = (offset - 1) * cols
	}
	gw := games[start:end]
	gwIDs := func() []string {
		r := []string{}
		for _, g := range gw {
			r = append(r, fmt.Sprintf("game_%v", g.Appid))
		}
		return r
	}()
	counter := 0
	lays := make([]*seker.Layer, len(app.Scene.LayersStack))
	copy(lays, app.Scene.LayersStack)
	toMove := []Mover{}
	for _, l := range lays {
		if l.Name == "root" {
			continue
		}
		b := false
		for _, n := range gwIDs {
			if l.Name == n {
				for _, item := range l.Items {
					toMove = append(toMove, Mover{*item, int32(-th * d)})
				}
				b = true
				counter++
				break
			}
		}
		if b {
			continue
		}
		// fmt.Println("r", l.Desc)
		app.Scene.RemoveLayer(l.Name)
	}
	// Shuffle(toMove)
	for _, m := range toMove {
		WG.Add(1)
		go func(mov Mover) {
			mov.Item.AnimateMove(0, mov.Y, 200)
			WG.Done()
		}(m)
		// m.Item.Move(0, m.Y)
	}
	c := 0
	r := rows
	gs := games[(offset+rows+1)*cols : (offset+rows+2)*cols]
	if d == -1 {
		r = -1
		gs = games[start : (offset+2)*cols]
	}
	for i, g := range gs {
		c = i % cols
		WG.Add(1)
		go func(row int, col int, game st.Game) {
			drawItem(app, game, col, row)
			WG.Done()
		}(r, c, g)
	}
	WG.Wait()
	LOCK.Unlock()
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if _, err := os.Stat("cache"); err != nil {
		os.Mkdir("cache", 0777)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	steam = st.Connect(viper.GetStringMapString("steam"))
	app := new(Application)
	os.Exit(app.run())
}
