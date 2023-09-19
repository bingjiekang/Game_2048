# Game_2048
 基于Golang的2048小游戏

参考：

[Ebitengine](https://github.com/hajimehoshi/ebiten)  
[同上](https://ebiten-zh.vercel.app/documents/mobile.html)

## 如何运行


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

1. 先编写背景及格子颜色

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
	
2. input.go


3. tile.go