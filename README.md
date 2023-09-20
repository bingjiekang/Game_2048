# Game_2048
 基于Golang的2048小游戏

参考：

[Ebitengine](https://github.com/hajimehoshi/ebiten)  
[同上](https://ebiten-zh.vercel.app/documents/mobile.html)

## 如何运行

### 方式1：配置好go环境、以及git后，克隆到本地，然后直接go run main.go

```
git clone https://github.com/bingjiekang/Game_2048.git
cd /Game_2048
go run main.go
```

### 方式2：下载压缩包到本地

```
1. 点击download下载到本地
2. 解压缩
3. 配置好go、c编译环境
4. cd /Game_2048
5. go run main.go
```

### 方式3：自己编写，代码如下（供参考）

## 如何实现

### 1.安装GO编辑器及C编辑器

```
参考：
https://ebiten-zh.vercel.app/documents/install.html
```

### 2.下载ebitengine包

```golang
// 创建一个文件夹
mkdir yourname
// 配置mod依赖
cd yourname
go mod init yourname
// 下载这个包
go get github.com/hajimehoshi/ebiten/v2
// 运行示例
go run -tags=example github.com/hajimehoshi/ebiten/v2/examples/rotate
// 出现一个旋转的图案则证明安装好了

// 导入包
go get "github.com/hajimehoshi/ebiten/v2/text"
```

### 3.代码编写

1. 编写main.go 主文件,用来启动游戏及打开窗口

	```golang
	// Game_2048/main.go
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
	
	```

2. 编写game.go 用来处理窗口大小，帧以及刷新频率相关

	```golang
	package twothousandandfortyeight
	// Game_2048/TwoThousandAndFortyEight/game.go
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
		input      *Input
		board      *Board
		boardImage *ebiten.Image
	}
	
	// 创建一个新游戏对象
	func NewGame() (*Game, error) {
		g := &Game{
			input: NewInput(),
		}
		var err error
		g.board, err = NewBoard(boardSize)
		if err != nil {
			return nil, err
		}
		return g, nil
	}
	
	// update控制更新当前帧率,每个tick都会被调用。tick是引擎更新的一个时间单位，默认为1/60s。tick的倒数我们一般称为帧，即游戏的更新频率。默认ebiten游戏是60帧，即每秒更新60次。
	func (g *Game) Update() error {
		g.input.Update()
		if err := g.board.Update(g.input); err != nil {
			return err
		}
		return nil
	}
	
	// 每帧（frame）调用。帧是渲染使用的一个时间单位，依赖显示器的刷新率。如果显示器的刷新率为60Hz，Draw将会每秒被调用60次
	func (g *Game) Draw(screen *ebiten.Image) {
		if g.boardImage == nil {
			g.boardImage = ebiten.NewImage(g.board.Size())
		}
		screen.Fill(BackgroundColor)
		g.board.Draw(g.boardImage)
		op := &ebiten.DrawImageOptions{}
		sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
		bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
		x := (sw - bw) / 2
		y := (sh - bh) / 2
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(g.boardImage, op)
	}
	
	// 该方法接收游戏窗口的尺寸作为参数，返回游戏的逻辑屏幕大小。
	func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
		return ScreenWidth, ScreenHeight
	}
	
	```

3. 编写背景及格子颜色

	```golang
	// Game_2048/TwoThousandAndFortyEight/colors.go
	// 格子颜色
	
	package twothousandandfortyeight
	
	/*
	游戏格子的颜色部分
	2048游戏极限为131072,要求新出的数都为4,不满足游戏要求,所以设置最大到65536即可
	*/
	import (
		"image/color"
	)
	
	var (
		// 背景及边框颜色
		BackgroundColor = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
		FrameColor      = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
	)
	
	// 对应数字字体的颜色2/4/8/16....,65536
	func TileColor(value int) color.Color {
		switch value {
		case 2, 4: // 字体颜色白色
			return color.RGBA{0x77, 0x6e, 0x65, 0xff}
		case 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536: // 字体颜色黑色
			return color.RGBA{0xf9, 0xf6, 0xf2, 0xff}
		}
		panic("数字颜色匹配不成功")
	}
	
	// 对应字体框颜色
	func TileBackgroundColor(value int) color.Color {
		switch value {
		case 0:
			return color.NRGBA{0xee, 0xe4, 0xda, 0x59}
		case 2:
			return color.RGBA{0xee, 0xe4, 0xda, 0xff}
		case 4:
			return color.RGBA{0xed, 0xe0, 0xc8, 0xff}
		case 8:
			return color.RGBA{0xf2, 0xb1, 0x79, 0xff}
		case 16:
			return color.RGBA{0xf5, 0x95, 0x63, 0xff}
		case 32:
			return color.RGBA{0xf6, 0x7c, 0x5f, 0xff}
		case 64:
			return color.RGBA{0xf6, 0x5e, 0x3b, 0xff}
		case 128:
			return color.RGBA{0xed, 0xcf, 0x72, 0xff}
		case 256:
			return color.RGBA{0xed, 0xcc, 0x61, 0xff}
		case 512:
			return color.RGBA{0xed, 0xc8, 0x50, 0xff}
		case 1024:
			return color.RGBA{0xed, 0xc5, 0x3f, 0xff}
		case 2048:
			return color.RGBA{0xed, 0xc2, 0x2e, 0xff}
		case 4096:
			return color.NRGBA{0xa3, 0x49, 0xa4, 0x7f}
		case 8192:
			return color.NRGBA{0xa3, 0x49, 0xa4, 0xb2}
		case 16384:
			return color.NRGBA{0xa3, 0x49, 0xa4, 0xcc}
		case 32768:
			return color.NRGBA{0xa3, 0x49, 0xa4, 0xe5}
		case 65536:
			return color.NRGBA{0xa3, 0x49, 0xa4, 0xff}
		}
		panic("数字框格匹配错误")
	}
	```
	
4. 对输入进行逻辑处理

	```golang
	package twothousandandfortyeight
	// Game_2048/TwoThousandAndFortyEight/input.go
	import (
		"github.com/hajimehoshi/ebiten/v2"
		"github.com/hajimehoshi/ebiten/v2/inpututil"
	)
	
	/* 
	对输入的上下左右进行获取和反馈处理
	 */
	
	type Dir int
	
	// 定义向上,向右,向下,向左 对应0-3
	const (
		DirUp Dir = iota
		DirRight
		DirDown
		DirLeft
	)
	
	type mouseState int
	
	const (
		mouseStateNone mouseState = iota
		mouseStatePressing
		mouseStateSettled
	)
	
	type touchState int
	
	const (
		touchStateNone touchState = iota
		touchStatePressing
		touchStateSettled
		touchStateInvalid
	)
	
	// 根据对应上下左右的数字返回对应移动字符
	func (d Dir) String() string {
		switch d {
		case DirUp:
			return "Up"
		case DirRight:
			return "Right"
		case DirDown:
			return "Down"
		case DirLeft:
			return "Left"
		}
		panic("无法到达")
	}
	
	// Vector为每个轴返回一个[-1，1]值。
	// 以左上角为原点坐标(0,0),x轴向右为正方向,y轴向下为正方向
	func (d Dir) Vector() (x, y int) {
		switch d {
		case DirUp:
			return 0, -1
		case DirRight:
			return 1, 0
		case DirDown:
			return 0, 1
		case DirLeft:
			return -1, 0
		}
		panic("无法到达")
	}
	
	// 获取当前输入的上下左右移动键
	type Input struct {
		mouseState    mouseState
		mouseInitPosX int
		mouseInitPosY int
		mouseDir      Dir
	
		touches       []ebiten.TouchID
		touchState    touchState
		touchID       ebiten.TouchID
		touchInitPosX int
		touchInitPosY int
		touchLastPosX int
		touchLastPosY int
		touchDir      Dir
	}
	
	// NewInput 是获取新的输入对象
	func NewInput() *Input {
		return &Input{}
	}
	
	func abs(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	
	func vecToDir(dx, dy int) (Dir, bool) {
		if abs(dx) < 4 && abs(dy) < 4 {
			return 0, false
		}
	
		if abs(dx) < abs(dy) {
			if dy < 0 {
				return DirUp, true
			}
			return DirDown, true
		}
	
		if dx < 0 {
			return DirLeft, true
		}
		return DirRight, true
	}
	
	// Update 更新当前的输入状态
	func (i *Input) Update() {
		switch i.mouseState {
		case mouseStateNone:
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				x, y := ebiten.CursorPosition()
				i.mouseInitPosX = x
				i.mouseInitPosY = y
				i.mouseState = mouseStatePressing
			}
		case mouseStatePressing:
			if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				x, y := ebiten.CursorPosition()
				dx := x - i.mouseInitPosX
				dy := y - i.mouseInitPosY
				d, ok := vecToDir(dx, dy)
				if !ok {
					i.mouseState = mouseStateNone
					break
				}
				i.mouseDir = d
				i.mouseState = mouseStateSettled
			}
		case mouseStateSettled:
			i.mouseState = mouseStateNone
		}
	
		i.touches = ebiten.AppendTouchIDs(i.touches[:0])
		switch i.touchState {
		case touchStateNone:
			if len(i.touches) == 1 {
				i.touchID = i.touches[0]
				x, y := ebiten.TouchPosition(i.touches[0])
				i.touchInitPosX = x
				i.touchInitPosY = y
				i.touchLastPosX = x
				i.touchLastPosY = y
				i.touchState = touchStatePressing
			}
		case touchStatePressing:
			if len(i.touches) >= 2 {
				break
			}
			if len(i.touches) == 1 {
				if i.touches[0] != i.touchID {
					i.touchState = touchStateInvalid
				} else {
					x, y := ebiten.TouchPosition(i.touches[0])
					i.touchLastPosX = x
					i.touchLastPosY = y
				}
				break
			}
			if len(i.touches) == 0 {
				dx := i.touchLastPosX - i.touchInitPosX
				dy := i.touchLastPosY - i.touchInitPosY
				d, ok := vecToDir(dx, dy)
				if !ok {
					i.touchState = touchStateNone
					break
				}
				i.touchDir = d
				i.touchState = touchStateSettled
			}
		case touchStateSettled:
			i.touchState = touchStateNone
		case touchStateInvalid:
			if len(i.touches) == 0 {
				i.touchState = touchStateNone
			}
		}
	}
	
	// Dir返回当前按下的方向。
	// 如果未按下方向键，Dir将返回false。
	func (i *Input) Dir() (Dir, bool) {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			return DirUp, true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			return DirLeft, true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			return DirRight, true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			return DirDown, true
		}
		if i.mouseState == mouseStateSettled {
			return i.mouseDir, true
		}
		if i.touchState == touchStateSettled {
			return i.touchDir, true
		}
		return 0, false
	}
	
	
	```


5. 对磁铁快进行合并及移动计算处理

	```golang
	package twothousandandfortyeight
	// Game_2048/TwoThousandAndFortyEight/tile.go
	import (
		"errors"
		"image/color"
		"log"
		"math/rand"
		"sort"
		"strconv"
	
		"github.com/hajimehoshi/ebiten/v2/text"
	
		"github.com/hajimehoshi/ebiten/v2"
		"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	
		"golang.org/x/image/font"
		"golang.org/x/image/font/opentype"
	)
	
	var (
		mplusSmallFont  font.Face // 小点的字体
		mplusNormalFont font.Face // 正常的字体
		mplusBigFont    font.Face // 大点的字体
	)
	
	const (
		maxMovingCount  = 5
		maxPoppingCount = 6
	)
	
	// 对应值和位置
	type TileData struct {
		value int
		x     int
		y     int
	}
	
	// 平铺数据和平铺动态
	type Tile struct {
		// 当前的位置及值
		current TileData
		// 下一个是空,如果还没移动到下一位
		next TileData
	
		movingCount       int
		startPoppingCount int
		poppingCount      int
	}
	
	// 初始化字体大小
	func init() {
		// 设置数字的字体
		tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}
	
		const dpi = 72
		// newface 返回一个新字体
		mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,                   // 字体大小
			DPI:     dpi,                  // 分辨率点数
			Hinting: font.HintingVertical, // 量化矢量字体的字形节点
		})
		if err != nil {
			log.Fatal(err)
		}
		// 正常大小的字体
		mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    32,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		// 大点的字体
		mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    48,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
	
	}
	
	// Pos返回当前位置
	func (t *Tile) Pos() (int, int) {
		return t.current.x, t.current.y
	}
	
	// 下一个移动到的位置
	func (t *Tile) NextPos() (int, int) {
		return t.next.x, t.next.y
	}
	
	// 当前位置的值
	func (t *Tile) Value() int {
		return t.current.value
	}
	
	// 下一个移动位置的值
	func (t *Tile) NextValue() int {
		return t.next.value
	}
	
	// NewTile 为新创建的一个移动对象
	func NewTile(value int, x, y int) *Tile {
		return &Tile{
			current: TileData{
				value: value,
				x:     x,
				y:     y,
			},
			startPoppingCount: maxPoppingCount,
		}
	}
	
	// 表示是否在设置动画(移动),返回一个bool
	func (t *Tile) IsMoving() bool {
		return 0 < t.movingCount
	}
	
	// 移动停止到的动画
	func (t *Tile) stopAnimation() {
		if 0 < t.movingCount {
			t.current = t.next
			t.next = TileData{}
		}
		// 然后初始化他们都为0
		t.movingCount = 0
		t.startPoppingCount = 0
		t.poppingCount = 0
	}
	
	// 获取当前坐标等于x,y的
	func tileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
		var result *Tile
		for t := range tiles {
			// 如果当前的x,y 不等于x或者y则进入下一个循环
			if t.current.x != x || t.current.y != y {
				continue
			}
			if result != nil {
				panic("无法到达")
			}
			result = t
		}
		return result
	}
	
	// 下一位置的x和y
	func currentOrNextTileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
		var result *Tile
		for t := range tiles {
			if 0 < t.movingCount {
				if t.next.x != x || t.next.y != y || t.next.value == 0 {
					continue
				}
			} else {
				if t.current.x != x || t.current.y != y {
					continue
				}
			}
			if result != nil {
				panic("无法到达")
			}
			result = t
		}
		return result
	}
	
	// MoveTiles 在给定的方格映射中移动
	// 如果有要移动的磁铁 则返回true
	// 调用时 所有磁铁都不移动
	func MoveTiles(tiles map[*Tile]struct{}, size int, dir Dir) bool {
		vx, vy := dir.Vector()
		tx := []int{}
		ty := []int{}
		for i := 0; i < size; i++ {
			tx = append(tx, i)
			ty = append(ty, i)
		}
		if vx > 0 {
			sort.Sort(sort.Reverse(sort.IntSlice(tx)))
		}
		if vy > 0 {
			sort.Sort(sort.Reverse(sort.IntSlice(ty)))
		}
	
		moved := false
		for _, j := range ty {
			for _, i := range tx {
				t := tileAt(tiles, i, j)
				if t == nil {
					continue
				}
				if t.next != (TileData{}) {
					panic("无法到达")
				}
				if t.IsMoving() {
					panic("无法到达")
				}
				//（ii，jj）是瓦片t的下一个位置。
				//（ii，jj）被更新，直到找到可合并瓦片为止
				// 或者瓦片t不能再移动了。
				ii := i
				jj := j
				for {
					ni := ii + vx
					nj := jj + vy
					if ni < 0 || ni >= size || nj < 0 || nj >= size {
						break
					}
					tt := currentOrNextTileAt(tiles, ni, nj)
					if tt == nil {
						ii = ni
						jj = nj
						moved = true
						continue
					}
					if t.current.value != tt.current.value {
						break
					}
					if 0 < tt.movingCount && tt.current.value != tt.next.value {
						// tt已与另一个磁贴合并。
						// 在此中断而不更新（ii，jj）。
						break
					}
					ii = ni
					jj = nj
					moved = true
					break
				}
	
				// next 是瓦片的下一个状态
				next := TileData{}
				next.value = t.current.value
				// 如果在下一个位置（ii，jj）有一个瓦片，
				// 这应该是可合并。让我们合并。
				if tt := currentOrNextTileAt(tiles, ii, jj); tt != t && tt != nil {
					next.value = t.current.value + tt.current.value
					tt.next.value = 0
					tt.next.x = ii
					tt.next.y = jj
					tt.movingCount = maxMovingCount
				}
				next.x = ii
				next.y = jj
				if t.current != next {
					t.next = next
					t.movingCount = maxMovingCount
				}
			}
		}
		if !moved {
			for t := range tiles {
				t.next = TileData{}
				t.movingCount = 0
			}
		}
		return moved
	}
	
	// 在屏幕上随机添加数
	func addRandomTile(tiles map[*Tile]struct{}, size int) error {
		cells := make([]bool, size*size)
		for t := range tiles {
			if t.IsMoving() {
				panic("无法实现")
			}
			i := t.current.x + t.current.y*size
			cells[i] = true
		}
		availableCells := []int{}
		for i, b := range cells {
			if b {
				continue
			}
			availableCells = append(availableCells, i)
		}
		if len(availableCells) == 0 {
			return errors.New("2048: there is no space to add a new tile")
		}
		c := availableCells[rand.Intn(len(availableCells))]
		v := 2
		if rand.Intn(10) == 0 {
			v = 4
		}
		x := c % size
		y := c / size
		t := NewTile(v, x, y)
		tiles[t] = struct{}{}
		return nil
	
	}
	
	// Update 更新当前格子的动画状态
	func (t *Tile) Update() error {
		switch {
		case 0 < t.movingCount:
			t.movingCount--
			if t.movingCount == 0 {
				if t.current.value != t.next.value && 0 < t.next.value {
					t.poppingCount = maxPoppingCount
				}
				t.current = t.next
				t.next = TileData{}
			}
		case 0 < t.startPoppingCount:
			t.startPoppingCount--
		case 0 < t.poppingCount:
			t.poppingCount--
		}
		return nil
	}
	
	func mean(a, b int, rate float64) int {
		return int(float64(a)*(1-rate) + float64(b)*rate)
	}
	
	func meanF(a, b float64, rate float64) float64 {
		return a*(1-rate) + b*rate
	}
	
	const (
		tileSize   = 80
		tileMargin = 4
	)
	
	var (
		tileImage = ebiten.NewImage(tileSize, tileSize)
	)
	
	func init() {
		tileImage.Fill(color.White)
	}
	
	// Draw将当前磁贴绘制到给定的板Image。
	func (t *Tile) Draw(boardImage *ebiten.Image) {
		i, j := t.current.x, t.current.y
		ni, nj := t.next.x, t.next.y
		v := t.current.value
		if v == 0 {
			return
		}
	
		op := &ebiten.DrawImageOptions{}
		x := i*tileSize + (i+1)*tileMargin
		y := j*tileSize + (j+1)*tileMargin
		nx := ni*tileSize + (ni+1)*tileMargin
		ny := nj*tileSize + (nj+1)*tileMargin
		switch {
		case 0 < t.movingCount:
			rate := 1 - float64(t.movingCount)/maxMovingCount
			x = mean(x, nx, rate)
			y = mean(y, ny, rate)
		case 0 < t.startPoppingCount:
			rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
			scale := meanF(0.0, 1.0, rate)
			op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
		case 0 < t.poppingCount:
			const maxScale = 1.2
			rate := 0.0
			if maxPoppingCount*2/3 <= t.poppingCount {
				// 0 to 1
				rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount)
			} else {
				// 1 to 0
				rate = float64(t.poppingCount) / float64(maxPoppingCount)
			}
			scale := meanF(1.0, maxScale, rate)
			op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
		}
		op.GeoM.Translate(float64(x), float64(y))
		op.ColorScale.ScaleWithColor(TileBackgroundColor(v))
		boardImage.DrawImage(tileImage, op)
		str := strconv.Itoa(v)
	
		f := mplusBigFont
		switch {
		case 3 < len(str):
			f = mplusSmallFont
		case 2 < len(str):
			f = mplusNormalFont
		}
	
		w := font.MeasureString(f, str).Floor()
		h := (f.Metrics().Ascent + f.Metrics().Descent).Floor()
		x += (tileSize - w) / 2
		y += (tileSize-h)/2 + f.Metrics().Ascent.Floor()
		text.Draw(boardImage, str, f, x, y, TileColor(v))
	}
	
	```

6. 对制定的磁吸块状态进行编写

	```golang
	package twothousandandfortyeight
	// Game_2048/TwoThousandAndFortyEight/board.go
	import (
		"errors"
	
		"github.com/hajimehoshi/ebiten/v2"
	)
	
	var taskTerminated = errors.New("2048: 任务已终止")
	
	// 自定义错误
	type task func() error
	
	// 游戏背景板
	type Board struct {
		size  int
		tiles map[*Tile]struct{}
		tasks []task
	}
	
	// NewBoard生成一个新Board，并给出大小
	func NewBoard(size int) (*Board, error) {
		b := &Board{
			size:  size,
			tiles: map[*Tile]struct{}{},
		}
		for i := 0; i < 2; i++ {
			if err := addRandomTile(b.tiles, b.size); err != nil {
				return nil, err
			}
		}
		return b, nil
	}
	
	func (b *Board) tileAt(x, y int) *Tile {
		return tileAt(b.tiles, x, y)
	}
	
	// update 更新board状态
	func (b *Board) Update(input *Input) error {
		for t := range b.tiles {
			if err := t.Update(); err != nil {
				return err
			}
		}
		if 0 < len(b.tasks) {
			t := b.tasks[0]
			if err := t(); err == taskTerminated {
				b.tasks = b.tasks[1:]
			} else if err != nil {
				return err
			}
			return nil
		}
		if dir, ok := input.Dir(); ok {
			if err := b.Move(dir); err != nil {
				return err
			}
		}
		return nil
	}
	
	// 将磁贴移动任务排入队列
	func (b *Board) Move(dir Dir) error {
		for t := range b.tiles {
			t.stopAnimation()
		}
		if !MoveTiles(b.tiles, b.size, dir) {
			return nil
		}
		b.tasks = append(b.tasks, func() error {
			for t := range b.tiles {
				if t.IsMoving() {
					return nil
				}
			}
			return taskTerminated
		})
		b.tasks = append(b.tasks, func() error {
			nextTiles := map[*Tile]struct{}{}
			for t := range b.tiles {
				if t.IsMoving() {
					panic("无法到达")
				}
				if t.next.value != 0 {
					panic("无法到达")
				}
				if t.current.value == 0 {
					continue
				}
				nextTiles[t] = struct{}{}
			}
			b.tiles = nextTiles
			if err := addRandomTile(b.tiles, b.size); err != nil {
				return err
			}
			return taskTerminated
		})
		return nil
	}
	
	// Size 返回这个board 的大小
	func (b *Board) Size() (int, int) {
		x := b.size*tileSize + (b.size+1)*tileMargin
		y := x
		return x, y
	}
	
	// Draw将板绘制到给定的板Image。
	func (b *Board) Draw(boardImage *ebiten.Image) {
		boardImage.Fill(FrameColor)
		for j := 0; j < b.size; j++ {
			for i := 0; i < b.size; i++ {
				v := 0
				op := &ebiten.DrawImageOptions{}
				x := i*tileSize + (i+1)*tileMargin
				y := j*tileSize + (j+1)*tileMargin
				op.GeoM.Translate(float64(x), float64(y))
				op.ColorScale.ScaleWithColor(TileBackgroundColor(v))
				boardImage.DrawImage(tileImage, op)
			}
		}
	
		animatingTiles := map[*Tile]struct{}{}
		nonAnimatingTiles := map[*Tile]struct{}{}
		for t := range b.tiles {
			if t.IsMoving() {
				animatingTiles[t] = struct{}{}
			} else {
				nonAnimatingTiles[t] = struct{}{}
			}
		}
		for t := range nonAnimatingTiles {
			t.Draw(boardImage)
		}
		for t := range animatingTiles {
			t.Draw(boardImage)
		}
	
	}
	```