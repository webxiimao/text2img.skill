package main

import "image/color"

// PageConfig 页面配置
type PageConfig struct {
	Width      int     // 页面宽度（像素）
	Height     int     // 页面高度（像素）
	PaddingX   int     // 水平内边距
	PaddingTop int     // 顶部内边距
	PaddingBot int     // 底部内边距
	FontSize   float64 // 字体大小
	LineHeight float64 // 行高倍数
	BgColor    color.Color
	TextColor  color.Color
	ShowPageNo bool   // 是否显示页码
	FontPath   string // 自定义字体路径
}

// Template 预设模板
type Template struct {
	Name   string
	Config PageConfig
}

// DefaultConfig 默认配置（小红书风格，竖屏）
func DefaultConfig() PageConfig {
	return PageConfig{
		Width:      1080,
		Height:     1440,
		PaddingX:   80,
		PaddingTop: 100,
		PaddingBot: 100,
		FontSize:   36,
		LineHeight: 1.8,
		BgColor:    color.White,
		TextColor:  color.RGBA{R: 51, G: 51, B: 51, A: 255},
		ShowPageNo: true,
	}
}

// 预设模板列表
var Templates = map[string]Template{
	"default": {
		Name:   "默认（白底黑字）",
		Config: DefaultConfig(),
	},
	"dark": {
		Name: "暗色",
		Config: PageConfig{
			Width: 1080, Height: 1440,
			PaddingX: 80, PaddingTop: 100, PaddingBot: 100,
			FontSize: 36, LineHeight: 1.8,
			BgColor:   color.RGBA{R: 30, G: 30, B: 30, A: 255},
			TextColor:  color.RGBA{R: 230, G: 230, B: 230, A: 255},
			ShowPageNo: true,
		},
	},
	"warm": {
		Name: "暖色纸张",
		Config: PageConfig{
			Width: 1080, Height: 1440,
			PaddingX: 80, PaddingTop: 100, PaddingBot: 100,
			FontSize: 36, LineHeight: 1.8,
			BgColor:   color.RGBA{R: 253, G: 245, B: 230, A: 255},
			TextColor:  color.RGBA{R: 60, G: 50, B: 40, A: 255},
			ShowPageNo: true,
		},
	},
}
