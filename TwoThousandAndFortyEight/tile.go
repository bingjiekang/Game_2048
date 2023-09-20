package twothousandandfortyeight

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
