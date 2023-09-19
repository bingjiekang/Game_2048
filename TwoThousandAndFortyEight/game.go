package twothousandandfortyeight

/*游戏开始部分
实现Game接口的Update/Draw/Layout函数,实现这个功能
*/

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// 游戏窗口的宽和高,以及边缘线
const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 4
)

type Game struct {
}

// update控制更新当前帧率,每个tick都会被调用。tick是引擎更新的一个时间单位，默认为1/60s。tick的倒数我们一般称为帧，即游戏的更新频率。默认ebiten游戏是60帧，即每秒更新60次。
func (g *Game) Update() error {
	return nil
}

// 每帧（frame）调用。帧是渲染使用的一个时间单位，依赖显示器的刷新率。如果显示器的刷新率为60Hz，Draw将会每秒被调用60次
func (g *Game) Draw(screen *ebiten.Image) {

}

// 该方法接收游戏窗口的尺寸作为参数，返回游戏的逻辑屏幕大小。
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
