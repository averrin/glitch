package main

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Drawable interface {
	Draw()
	Move(int32, int32)
	AnimateMove(int32, int32, int32)
	MoveTo(int32, int32)
	IsChanged() bool
	Clear()
	SetScale(float64)
	GetScale() float64
	GetRect() *sdl.Rect
	SetParentSurface(*sdl.Surface)
	Destroy()
}

type Rect struct {
	ParentSurface *sdl.Surface
	Rect          *sdl.Rect
	LastRects     []*sdl.Rect
	Color         uint32
	Scale         float64
	LastScale     float64
	Changed       bool
}

type Image struct {
	Rect
	Path  string
	Image *sdl.Surface
}

func (item *Rect) SetParentSurface(s *sdl.Surface) {
	item.ParentSurface = s
}

func StripLine(line string, w int32) string {
	lw, _, _ := font.SizeUTF8(line)
	for int32(lw) > int32(w)-16 {
		line = strings.TrimRight(line[:len(line)-4], " -") + "…"
		lw, _, _ = font.SizeUTF8(line)
	}
	return line
}

func NewImage(rect *sdl.Rect, path string, alt string) Image {
	item := new(Image)
	item.Rect = NewRect(rect, 0xff000000)
	item.Path = path
	image, _ := sdl.LoadBMP(item.Path)
	item.Image = image
	if item.Image == nil {
		amask := uint32(0xff000000)
		rmask := uint32(0x00ff0000)
		gmask := uint32(0x0000ff00)
		bmask := uint32(0x000000ff)
		s, _ := sdl.CreateRGBSurface(sdl.SWSURFACE, rect.W, rect.H, 32, rmask, gmask, bmask, amask)
		item.Image = s
		item.Image.FillRect(&sdl.Rect{0, 0, rect.W, rect.H}, 0xff000000)
		alt = StripLine(alt, rect.W)
		lw, _, _ := font.SizeUTF8(alt)
		title := NewText(&sdl.Rect{int32(rect.W/2) - int32(lw/2), int32(rect.H/2) - int32(font.Height()/2), int32(lw), int32(font.Height())}, alt, 0xfff0f0f0)
		title.SetParentSurface(item.Image)
		title.Draw()
	}
	return *item
}

func NewRect(rect *sdl.Rect, color uint32) Rect {
	item := new(Rect)
	item.Rect = rect
	item.Color = color
	item.Scale = 1
	item.LastRects = []*sdl.Rect{item.Rect}
	item.Changed = true
	return *item
}

func (item *Image) Draw() {
	item.Clear()
	s := item.ParentSurface
	r := item.GetRect()
	// log.Println(r.X/r.W, r.Y/r.H)
	item.Image.BlitScaled(
		&sdl.Rect{0, 0, item.Image.W, item.Image.H},
		s,
		&sdl.Rect{r.X, r.Y, int32(float64(r.W) * item.Scale), int32(float64(r.H) * item.Scale)},
	)
	item.Changed = false
	item.LastRects = append(item.LastRects, item.GetRect())
}

func (item *Image) Destroy() {
	item.Image.Free()
}

func (item *Rect) Destroy() {
}

func (item *Rect) Draw() {
	s := item.ParentSurface
	s.FillRect(item.Rect, item.Color)
	item.Changed = false
	item.LastRects = append(item.LastRects, item.GetRect())
}

func (item *Rect) GetLastRects() []*sdl.Rect {
	return item.LastRects
}

func (item *Rect) GetRect() *sdl.Rect {
	return item.Rect
}

func (item *Rect) SetRect(rect *sdl.Rect) {
	item.LastRects = append(item.LastRects, item.GetRect())
	item.Rect = rect
	item.Changed = true
}

func (item *Rect) IsChanged() bool {
	return item.Changed
}

func (item *Rect) Clear() {
	s := item.ParentSurface
	for _, r := range item.LastRects {
		s.FillRect(r, 0x00000000)
	}
	item.LastRects = []*sdl.Rect{}
	// lr := sdl.Rect{r.X, r.Y, int32(float64(r.W) * item.LastScale), int32(float64(r.H) * item.LastScale)}
}

func (item *Rect) SetScale(scale float64) {
	item.LastScale = item.Scale
	item.Scale = scale
	item.Changed = true
}

func (item *Rect) GetScale() float64 {
	return item.Scale
}

func (item *Rect) AnimateMove(x int32, y int32, duration int32) {
	duration = duration / 2
	dx := int32(float32(x) / float32(duration))
	dy := int32(float32(y) / float32(duration))
	ey := item.Rect.Y + y
	ex := item.Rect.X + x
	for dx != 0 || dy != 0 {
		if math.Abs(float64(item.Rect.Y-ey)) < math.Abs(float64(dy)) {
			dy = item.Rect.Y - ey
			if math.Abs(float64(dy)) < 1 {
				dy = 0
			}
		}
		if math.Abs(float64(item.Rect.X-ex)) < math.Abs(float64(dx)) {
			dx = item.Rect.X - ex
			if math.Abs(float64(dx)) < 1 {
				dx = 0
			}
		}

		item.Move(dx, dy)
		time.Sleep(1 * time.Millisecond)
	}
}

func (item *Rect) Move(x int32, y int32) {
	item.LastRects = append(item.LastRects, &sdl.Rect{item.Rect.X, item.Rect.Y, item.Rect.W, item.Rect.H})
	item.Rect.X += x
	item.Rect.Y += y
	item.Changed = true
}

func (item *Rect) MoveTo(x int32, y int32) {
	item.LastRects = append(item.LastRects, &sdl.Rect{item.Rect.X, item.Rect.Y, item.Rect.W, item.Rect.H})
	item.Rect.X = x
	item.Rect.Y = y
	item.Changed = true
}

type Text struct {
	Rect
	Text string
}

func NewText(rect *sdl.Rect, text string, color uint32) Text {
	item := new(Text)
	item.Rect = NewRect(rect, color)
	item.Text = text
	return *item
}

func (item *Text) Draw() {
	s := item.ParentSurface
	message, err := font.RenderUTF8_Blended(item.Text, sdl.Color{250, 250, 250, 1})
	if err != nil {
		log.Fatal(err)
	}
	defer message.Free()
	srcRect := sdl.Rect{}
	message.GetClipRect(&srcRect)
	message.Blit(&srcRect, s, item.GetRect())
	item.Changed = false
	item.LastRects = append(item.LastRects, item.Rect.Rect)
}
