package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	rand.Seed(time.Now().Unix())
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(1920, 1080, sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	f, err := ttf.OpenFont("res/fonts/BabelStoneHan.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	s := newScene(r, w, f)
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(events)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

type scene struct {
	font     *ttf.Font
	renderer *sdl.Renderer
	window   *sdl.Window
	lines    []*line
}

func newScene(r *sdl.Renderer, w *sdl.Window, f *ttf.Font) *scene {
	mlines := []*line{}
	width, _ := w.GetSize()
	count := int(width / 40)
	for i := 0; i < count; i++ {
		mlines = append(mlines, newLine(int32(i*40), r, f))
	}

	return &scene{
		font:     f,
		renderer: r,
		window:   w,
		lines:    mlines,
	}
}

func (s *scene) run(events <-chan sdl.Event) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(5 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if err := s.paint(); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyUpEvent:
		if event.(*sdl.KeyUpEvent).Keysym.Sym == 27 {
			return true
		}
	}
	return false
}

func (s *scene) update() {
	for _, l := range s.lines {
		l.update(s.window)
	}
}

func (s *scene) paint() error {
	s.renderer.Clear()
	for _, l := range s.lines {
		l.paint(s.window)
	}
	s.renderer.Present()
	return nil
}

func (s *scene) destroy() {
}
