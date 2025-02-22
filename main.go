package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/aldernero/gaul"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	maxSpeed = 10
	maxTime  = 30
)

//go:embed assets/*.png
var assets embed.FS

var (
	W int
	H int
) // window size

type state struct {
	pos gaul.Vec2
	vel gaul.Vec2
}

type game struct {
	img          *ebiten.Image
	images       []*ebiten.Image
	started      time.Time
	screenWidth  float64
	screenHeight float64
	imgWidth     float64
	imgHeight    float64
	state
}

// Setup initializes the game state.
func (g *game) Setup() {
	g.img = g.images[rand.IntN(len(g.images))]
	g.imgWidth, g.imgHeight = float64(g.img.Bounds().Dx()), float64(g.img.Bounds().Dy())
	g.pos = gaul.Vec2{X: g.screenWidth / 2, Y: g.screenHeight / 2}
	xyAngle := gaul.Tau * rand.Float64()
	g.vel = gaul.Vec2{X: maxSpeed * math.Cos(xyAngle), Y: maxSpeed * math.Sin(xyAngle)}
	g.started = time.Now()
}

func (g *game) Update() error {
	if time.Now().After(g.started.Add(maxTime * time.Second)) {
		g.Setup()
		return nil
	}
	g.pos = g.pos.Add(g.vel)
	if g.pos.X <= 0 || g.pos.X+g.imgWidth >= g.screenWidth {
		g.vel.X = -g.vel.X
	}
	if g.pos.Y <= 0 || g.pos.Y+g.imgHeight >= g.screenHeight {
		g.vel.Y = -g.vel.Y
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.pos.X, g.pos.Y)
	screen.DrawImage(g.img, op)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return W, H
}

func main() {
	// set window size to monitor size
	monitor := ebiten.Monitor()
	if monitor == nil {
		log.Fatal("no monitor found")
	}
	W, H = monitor.Size()
	ebiten.SetWindowSize(W, H)
	ebiten.SetFullscreen(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetVsyncEnabled(true)

	game := game{
		screenWidth:  float64(W),
		screenHeight: float64(H),
	}
	// open images
	imgs, err := assets.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
	}
	// load images
	game.images = make([]*ebiten.Image, len(imgs))
	for i, img := range imgs {
		imgBytes, err := assets.ReadFile(fmt.Sprintf("assets/%s", img.Name()))
		if err != nil {
			log.Fatal(err)
		}
		game.images[i], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(imgBytes))
		if err != nil {
			log.Fatal(err)
		}
	}

	game.Setup()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
