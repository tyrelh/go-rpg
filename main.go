package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	X                int
	Y                int
	PlayerImage      *ebiten.Image
	PlayerX, PlayerY float64
}

func (g *Game) Update() error {
	// react to key presses

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.PlayerX += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.PlayerX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.PlayerY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.PlayerY += 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	// draw player image

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(g.PlayerX, g.PlayerY)
	screen.DrawImage(
		g.PlayerImage.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&options,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.X, g.Y
	// return ebiten.WindowSize()
}

func main() {
	game := &Game{
		X:       240,
		Y:       160,
		PlayerX: 100,
		PlayerY: 100,
	}
	ebiten.SetWindowSize(game.X*4, game.Y*4)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello, World!")

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/raccoon-kit.png")
	if err != nil {
		log.Fatalf("Failed to load player image: %v", err)
	}
	game.PlayerImage = playerImage

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
