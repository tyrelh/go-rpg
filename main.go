package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Image *ebiten.Image
	X, Y  float64
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	HealAmount uint
}

type Player struct {
	*Sprite
	Health uint
}

type Game struct {
	X            int
	Y            int
	Player       *Player
	enemies      []*Enemy
	potions      []*Potion
	Tilemap      *TilemapJSON
	TilemapImage *ebiten.Image
}

func (g *Game) Update() error {
	// react to key presses

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Player.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Player.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Player.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Player.Y += 1
	}

	speed := 0.5
	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.Player.X {
				sprite.X += speed
			}
			if sprite.X > g.Player.X {
				sprite.X -= speed
			}
			if sprite.Y < g.Player.Y {
				sprite.Y += speed
			}
			if sprite.Y > g.Player.Y {
				sprite.Y -= speed
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	// draw player image

	opts := ebiten.DrawImageOptions{}

	// loop over tilemap layers
	for _, layer := range g.Tilemap.Layers {
		// loop over tiles in layer
		for i, id := range layer.Data {
			x := i % layer.Width
			y := i / layer.Width
			x *= 16
			y *= 16
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22
			srcX *= 16
			srcY *= 16
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(
				g.TilemapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
	}

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(g.Player.X, g.Player.Y)
	screen.DrawImage(
		g.Player.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&options,
	)

	for _, sprite := range g.enemies {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(
			sprite.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&options,
		)
	}

	for _, potion := range g.potions {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(potion.X, potion.Y)
		screen.DrawImage(
			potion.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&options,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.X, g.Y
	// return ebiten.WindowSize()
}

func main() {
	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/raccoon-kit.png")
	if err != nil {
		log.Fatalf("Failed to load player image: %v", err)
	}

	skeletonImage, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatalf("Failed to load player image: %v", err)
	}

	potionImage, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatalf("Failed to load potion image: %v", err)
	}

	tilemapImage, _, err := ebitenutil.NewImageFromFile("assets/images/Ninja Adventure - Asset Pack/Backgrounds/Tilesets/TilesetFloor.png")
	if err != nil {
		log.Fatalf("Failed to load tilemap image: %v", err)
	}

	tilemap, err := NewTilemapJSON("assets/maps/main.json")
	if err != nil {
		log.Fatalf("Failed to load tilemap: %v", err)
	}

	game := &Game{
		X: 240,
		Y: 160,
		Player: &Player{
			Sprite: &Sprite{
				Image: playerImage,
				X:     100,
				Y:     100,
			},
			Health: 3,
		},
		enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Image: skeletonImage,
					X:     50,
					Y:     50,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &Sprite{
					Image: skeletonImage,
					X:     100,
					Y:     100,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*Potion{
			{
				Sprite: &Sprite{
					Image: potionImage,
					X:     100,
					Y:     70,
				},
				HealAmount: 10,
			},
		},
		Tilemap:      tilemap,
		TilemapImage: tilemapImage,
	}

	ebiten.SetWindowSize(game.X*3, game.Y*3)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
