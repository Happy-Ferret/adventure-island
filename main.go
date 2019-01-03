package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/kyeett/adventure-island/resources"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 160
	screenHeight = 160
)

const (
	tileSize = 16
	tileXNum = 2
)

const xNum = screenWidth / tileSize

var world Map

var (
	tilesImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.All_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

var p Player

var (
	Up    = Coord{0, -1}
	Left  = Coord{-1, 0}
	Down  = Coord{0, 1}
	Right = Coord{1, 0}
)

func randomWalk() {

	var c Coord

	if time.Now().Nanosecond()%10 == 0 {
		switch rand.Intn(4) {
		case 0:
			c = p.PrepareMove(Up)
		case 1:
			c = p.PrepareMove(Left)
		case 2:
			c = p.PrepareMove(Down)
		case 3:
			c = p.PrepareMove(Right)
		}

		if world.ValidTarget(c) == true {
			p.MoveTo(c)
		}
	}

}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	var c Coord

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Game terminated by player")
	}

	randomWalk()

	if leftPressed() || rightPressed() || upPressed() || downPressed() {
		switch {
		case upPressed():
			c = p.PrepareMove(Up)
		case leftPressed():
			c = p.PrepareMove(Left)
		case downPressed():
			c = p.PrepareMove(Down)
		case rightPressed():
			c = p.PrepareMove(Right)
		}

		if world.ValidTarget(c) == true {
			p.MoveTo(c)
		}

	}

	world.CheckCollisions()

	world.Draw(screen)
	for _, o := range world.entities {
		o.Draw(screen)
	}

	// p.Draw(screen)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", world.GetScore()))

	return nil
}

func main() {

	world = NewMap()
	p = world.AddPlayer()

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Tiles (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
