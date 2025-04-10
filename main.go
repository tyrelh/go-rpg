package main

import (
	"go-rpg/entities"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	X            int
	Y            int
	Player       *entities.Player
	enemies      []*entities.Enemy
	potions      []*entities.Potion
	Tilemap      *TilemapJSON
	TilemapImage *ebiten.Image
}

func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		log.Println("Window is being closed...")
		// Perform any cleanup operations here
		os.Exit(0)
	}
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
			// get the position of the tile
			x := i % layer.Width
			y := i / layer.Width
			// convert to pixels
			x *= 16
			y *= 16
			// get the position on the image where the tile id is
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22
			// convert to pixels
			srcX *= 16
			srcY *= 16
			// set the x, y for drawing tile
			opts.GeoM.Translate(float64(x), float64(y))
			// draw the tile
			screen.DrawImage(
				g.TilemapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			// reset the geo m for next tile
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
		Player: &entities.Player{
			Sprite: &entities.Sprite{
				Image: playerImage,
				X:     100,
				Y:     100,
			},
			Health: 3,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Image: skeletonImage,
					X:     50,
					Y:     50,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Image: skeletonImage,
					X:     100,
					Y:     100,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exiting...")
		os.Exit(0)
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
