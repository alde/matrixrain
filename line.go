package main

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type line struct {
	render *sdl.Renderer
	font   *ttf.Font
	stop   <-chan int

	x       int32
	y       int32
	speed   int32
	symbols []*symbol
}

func newLine(x int32, r *sdl.Renderer, f *ttf.Font, stop <-chan int) *line {
	l := &line{
		render: r,
		font:   f,
		x:      x,
		stop:   stop,
	}

	l.reset()

	return l
}

func (l *line) update(w *sdl.Window) {
	_, height := w.GetSize()
	l.y += l.speed
	if l.y > int32(height)+int32(40*len(l.symbols)) {
		l.reset()
	}
}

func (l *line) paint(w *sdl.Window) error {
	rect := sdl.Rect{X: int32(l.x), Y: int32(l.y), W: 40, H: 500}

	for i, sym := range l.symbols {
		sym.paint(l.x, l.y-(40*int32(i)), l.render, l.font)
	}

	if err := l.render.Copy(l.render.GetRenderTarget(), nil, &rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}
	return nil
}

func (l *line) reset() {
	l.y = -rand.Int31n(1000)
	l.speed = rand.Int31n(20) + 20
	l.symbols = []*symbol{newSymbol(&sdl.Color{R: 120, G: 255, B: 120, A: 255}, l.stop)}

	symbolCount := rand.Intn(20) + 20
	for i := 0; i < symbolCount; i++ {
		c := &sdl.Color{R: 0, G: 255, B: 50, A: 255}
		l.symbols = append(l.symbols, newSymbol(c, l.stop))
	}
}
