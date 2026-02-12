package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var (
		inputFile  string
		outputDir  string
		template   string
		fontPath   string
		fontSize   float64
		lineHeight float64
		listTpl    bool
	)

	flag.StringVar(&inputFile, "i", "", "输入文本文件路径（不指定则从stdin读取）")
	flag.StringVar(&outputDir, "o", "./output", "输出目录")
	flag.StringVar(&template, "t", "default", "模板名称")
	flag.StringVar(&fontPath, "font", "", "自定义字体文件路径")
	flag.Float64Var(&fontSize, "size", 0, "字体大小（覆盖模板设置）")
	flag.Float64Var(&lineHeight, "lh", 0, "行高倍数（覆盖模板设置）")
	flag.BoolVar(&listTpl, "list", false, "列出所有可用模板")
	flag.Parse()

	if listTpl {
		fmt.Println("可用模板:")
		for k, v := range Templates {
			fmt.Printf("  %-10s %s\n", k, v.Name)
		}
		return
	}

	// 读取输入文本
	text, err := readInput(inputFile, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(text) == "" {
		fmt.Fprintln(os.Stderr, "错误: 输入内容为空")
		os.Exit(1)
	}

	// 加载配置
	tpl, ok := Templates[template]
	if !ok {
		fmt.Fprintf(os.Stderr, "错误: 未知模板 %q，使用 -list 查看可用模板\n", template)
		os.Exit(1)
	}
	cfg := tpl.Config
	if fontPath != "" {
		cfg.FontPath = fontPath
	}
	if fontSize > 0 {
		cfg.FontSize = fontSize
	}
	if lineHeight > 0 {
		cfg.LineHeight = lineHeight
	}

	// 加载字体
	face, err := LoadFont(cfg.FontPath, cfg.FontSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 排版分页
	pages := Layout(text, cfg, face)
	fmt.Printf("共分 %d 页\n", len(pages))

	// 渲染输出
	paths, err := Render(pages, cfg, face, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	for _, p := range paths {
		fmt.Println(p)
	}
	fmt.Printf("完成！共生成 %d 张图片，保存在 %s\n", len(paths), outputDir)
}

func readInput(filePath string, args []string) (string, error) {
	// 优先从文件读取
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("读取文件失败: %w", err)
		}
		return string(data), nil
	}

	// 其次从剩余参数读取
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// 最后从stdin读取
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("读取stdin失败: %w", err)
		}
		return string(data), nil
	}

	return "", fmt.Errorf("请指定输入文件(-i)或通过管道传入文本")
}
