package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Render 将所有页面渲染为图片并保存
func Render(pages []Page, cfg PageConfig, face font.Face, outputDir string) ([]string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %w", err)
	}

	var paths []string
	for _, page := range pages {
		img := renderPage(page, cfg, face)
		filename := fmt.Sprintf("page-%02d.png", page.Index)
		outPath := filepath.Join(outputDir, filename)
		if err := savePNG(img, outPath); err != nil {
			return nil, fmt.Errorf("保存第%d页失败: %w", page.Index, err)
		}
		paths = append(paths, outPath)
	}
	return paths, nil
}

func renderPage(page Page, cfg PageConfig, face font.Face) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))

	// 填充背景
	fillRect(img, cfg.BgColor)

	metrics := face.Metrics()
	ascent := metrics.Ascent.Ceil()

	// 计算可用内容区域高度
	contentHeight := cfg.Height - cfg.PaddingTop - cfg.PaddingBot
	if cfg.ShowPageNo && page.Total > 1 {
		contentHeight -= 50
	}

	// 动态计算行高，让文字均匀铺满页面
	numLines := len(page.Lines)
	lineHeightPx := int(cfg.FontSize * cfg.LineHeight) // 基础行高
	if numLines > 1 {
		// 用可用高度除以行数，得到实际行距
		dynamicLH := contentHeight / numLines
		// 不要让行距小于基础行高（避免文字重叠）
		if dynamicLH > lineHeightPx {
			lineHeightPx = dynamicLH
		}
	}

	// 计算起始Y偏移，使文字在内容区域内垂直居中
	totalTextHeight := (numLines - 1) * lineHeightPx
	startY := cfg.PaddingTop + (contentHeight-totalTextHeight)/2 + ascent

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(cfg.TextColor),
		Face: face,
	}

	for i, line := range page.Lines {
		x := cfg.PaddingX
		y := startY + i*lineHeightPx
		d.Dot = fixed.P(x, y)
		d.DrawString(line)
	}

	// 绘制页码
	if cfg.ShowPageNo && page.Total > 1 {
		pageNoStr := fmt.Sprintf("%d / %d", page.Index, page.Total)
		pageNoColor := color.RGBA{R: 180, G: 180, B: 180, A: 255}
		d.Src = image.NewUniform(pageNoColor)
		adv := font.MeasureString(face, pageNoStr)
		x := (cfg.Width - adv.Ceil()) / 2
		y := cfg.Height - cfg.PaddingBot/2
		d.Dot = fixed.P(x, y)
		d.DrawString(pageNoStr)
	}

	return img
}

func fillRect(img *image.RGBA, c color.Color) {
	draw.Draw(img, img.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
}

func savePNG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
