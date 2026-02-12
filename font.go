package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// LoadFont 加载字体，优先自定义路径，其次系统字体
func LoadFont(customPath string, size float64) (font.Face, error) {
	if customPath != "" {
		return loadFontFromFile(customPath, size)
	}
	return loadSystemFont(size)
}

func loadFontFromFile(path string, size float64) (font.Face, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取字体文件失败: %w", err)
	}
	return parseFontData(data, size)
}

func parseFontData(data []byte, size float64) (font.Face, error) {
	// 尝试作为单个字体解析
	f, err := opentype.Parse(data)
	if err == nil {
		return opentype.NewFace(f, &opentype.FaceOptions{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	}

	// 尝试作为字体集合解析
	col, err := opentype.ParseCollection(data)
	if err != nil {
		return nil, fmt.Errorf("无法解析字体文件: %w", err)
	}
	ft, err := col.Font(0)
	if err != nil {
		return nil, fmt.Errorf("无法从字体集合中提取字体: %w", err)
	}
	return opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// loadSystemFont 按平台搜索系统中文字体
func loadSystemFont(size float64) (font.Face, error) {
	candidates := systemFontPaths()
	for _, p := range candidates {
		matches, _ := filepath.Glob(p)
		for _, m := range matches {
			face, err := loadFontFromFile(m, size)
			if err == nil {
				return face, nil
			}
		}
	}
	return nil, fmt.Errorf("未找到可用的中文字体，请通过 --font 指定字体文件路径")
}

func systemFontPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/STHeiti Light.ttc",
			"/System/Library/Fonts/Supplemental/Songti.ttc",
			"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
			"/Library/Fonts/Arial Unicode.ttf",
		}
	case "windows":
		return []string{
			`C:\Windows\Fonts\msyh.ttc`,
			`C:\Windows\Fonts\msyhbd.ttc`,
			`C:\Windows\Fonts\simsun.ttc`,
			`C:\Windows\Fonts\simhei.ttf`,
		}
	case "linux":
		return []string{
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/google-noto-cjk/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		}
	default:
		return nil
	}
}
