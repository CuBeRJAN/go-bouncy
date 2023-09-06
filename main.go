package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const WIN_H int32 = 720
const WIN_W int32 = 1280
const TICK_INTERVAL uint32 = 5
const IMAGE_PATH string = "image.png"

var next_time uint32

type pos_t struct {
	x, y, dx, dy int32
}

func time_left() uint32 {
	var now uint32 = sdl.GetTicks()
	if next_time <= now {
		return 0
	}
	return next_time - now
}

func cycle_color(tex *sdl.Texture) {
	var colors [][]uint8 = [][]uint8{{255, 50, 255}, {50, 50, 255}, {255, 50, 50}, {50, 255, 255}, {255, 255, 255}}
	var rand_col []uint8 = colors[rand.Intn(len(colors))]
	tex.SetColorMod(rand_col[0], rand_col[1], rand_col[2])
}

func move_tex(tex *sdl.Texture, dst *sdl.Rect, dx *int32, dy *int32) {
	for true {
		dst.X += *dx
		dst.Y += *dy
		if dst.X+180 > WIN_W || dst.X < 0 {
			*dx *= -1
			cycle_color(tex)
		}
		if dst.Y+134 > WIN_H || dst.Y+36 < 0 {
			*dy *= -1
			cycle_color(tex)
		}
		sdl.Delay(time_left())
		next_time += TICK_INTERVAL
	}
}

func main() {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var texture *sdl.Texture
	var image *sdl.Surface
	var src, dst sdl.Rect

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIN_W, WIN_H, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	image, err = img.Load(IMAGE_PATH)
	if err != nil {
		panic(err)
	}
	defer image.Free()

	texture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	src = sdl.Rect{0, 0, 225, 225}
	dst = sdl.Rect{100, 50, 180, 180}

	var dx, dy int32
	dx = 1
	dy = 1
	go move_tex(texture, &dst, &dx, &dy)

	var texr sdl.Rect
	texr.X = WIN_W / 2
	texr.Y = WIN_H / 2

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.Copy(texture, &src, &dst)
		if err != nil {
			panic(err)
		}
		renderer.Present()
	}
}
