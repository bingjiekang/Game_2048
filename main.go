package main

import (
	twothousandandfortyeight "Game_2048/TwoThousandAndFortyEight"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := twothousandandfortyeight.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(twothousandandfortyeight.ScreenWidth, twothousandandfortyeight.ScreenHeight)
	ebiten.SetWindowTitle("2048")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
