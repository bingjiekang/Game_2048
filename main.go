package main

import (
	twothousandandfortyeight "Game_2048/TwoThousandAndFortyEight"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 创建一个游戏对象
	game, err := twothousandandfortyeight.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	// 游戏窗口大小
	ebiten.SetWindowSize(twothousandandfortyeight.ScreenWidth, twothousandandfortyeight.ScreenHeight)
	// 游戏标题
	ebiten.SetWindowTitle("2048")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
