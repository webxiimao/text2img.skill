package main

import (
	"strings"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Page 一页的内容
type Page struct {
	Lines []string
	Index int // 页码（从1开始）
	Total int // 总页数
}

// Layout 执行排版：换行 + 分页
func Layout(text string, cfg PageConfig, face font.Face) []Page {
	maxWidth := cfg.Width - cfg.PaddingX*2
	lineHeightPx := int(cfg.FontSize * cfg.LineHeight)
	contentHeight := cfg.Height - cfg.PaddingTop - cfg.PaddingBot
	if cfg.ShowPageNo {
		contentHeight -= 50 // 页码区域
	}
	maxLinesPerPage := contentHeight / lineHeightPx

	// 按段落分割，逐段换行
	paragraphs := strings.Split(text, "\n")
	var allLines []string
	for _, para := range paragraphs {
		para = strings.TrimRight(para, " \t\r")
		if para == "" {
			allLines = append(allLines, "")
			continue
		}
		wrapped := wrapLine(para, face, maxWidth)
		allLines = append(allLines, wrapped...)
	}

	// 分页
	return paginate(allLines, maxLinesPerPage)
}

// wrapLine 将一行文字按宽度换行
func wrapLine(text string, face font.Face, maxWidth int) []string {
	if text == "" {
		return []string{""}
	}

	runes := []rune(text)
	maxW := fixed.I(maxWidth)
	var lines []string
	start := 0

	for start < len(runes) {
		end := start
		var lastBreakable int = -1
		var w fixed.Int26_6

		for end < len(runes) {
			adv, ok := face.GlyphAdvance(runes[end])
			if !ok {
				adv = face.Metrics().Height / 2
			}
			w += adv

			if w > maxW && end > start {
				// 超出宽度，在合适位置断行
				if lastBreakable > start {
					lines = append(lines, string(runes[start:lastBreakable]))
					start = lastBreakable
				} else {
					lines = append(lines, string(runes[start:end]))
					start = end
				}
				break
			}

			// 记录可断行位置
			if isBreakable(runes, end) {
				lastBreakable = end + 1
			}
			end++
		}

		if end >= len(runes) {
			lines = append(lines, string(runes[start:]))
			break
		}
	}

	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// isBreakable 判断是否可以在此位置断行
func isBreakable(runes []rune, i int) bool {
	r := runes[i]
	// 空格后可断行
	if unicode.IsSpace(r) {
		return true
	}
	// CJK字符后可断行（中日韩统一表意文字）
	if isCJK(r) {
		// 但下一个字符如果是不能放在行首的标点，则不断
		if i+1 < len(runes) && isNoBreakBefore(runes[i+1]) {
			return false
		}
		return true
	}
	// 标点符号后可断行（但某些标点不能放在行尾）
	if unicode.IsPunct(r) && !isNoBreakAfter(r) {
		return true
	}
	return false
}

func isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK基本
		(r >= 0x3400 && r <= 0x4DBF) || // CJK扩展A
		(r >= 0x3000 && r <= 0x303F) || // CJK标点
		(r >= 0xFF00 && r <= 0xFFEF) // 全角字符
}

// 不能放在行首的字符
func isNoBreakBefore(r rune) bool {
	return strings.ContainsRune("，。！？、；：）》」』】〉〕）]}>,.!?;:)", r)
}

// 不能放在行尾的字符
func isNoBreakAfter(r rune) bool {
	return strings.ContainsRune("（《「『【〈〔（[{<", r)
}

// paragraph 一个段落（可能包含多行换行后的文字）
type paragraph struct {
	lines []string
}

// paginate 将行列表均匀分成多页，优先在段落边界分页
func paginate(lines []string, maxLinesPerPage int) []Page {
	if maxLinesPerPage <= 0 {
		maxLinesPerPage = 1
	}

	lines = trimEmptyLines(lines)
	if len(lines) == 0 {
		return nil
	}

	paras := splitParagraphs(lines)
	totalLines := countDisplayLines(paras)

	numPages := (totalLines + maxLinesPerPage - 1) / maxLinesPerPage
	if numPages <= 1 {
		return []Page{{Lines: buildPageLines(paras), Index: 1, Total: 1}}
	}

	// 计算每个段落的累积行数（含段落间空行）
	cumLines := make([]int, len(paras)+1) // cumLines[i] = 前i个段落的总显示行数
	for i, p := range paras {
		cumLines[i+1] = cumLines[i] + len(p.lines)
		if i > 0 {
			cumLines[i+1]++ // 段落间空行
		}
	}

	// 用DP找最优分割点，使各页行数差异最小
	// split[p] = 第p页结束后的段落索引
	bestSplit := findOptimalSplit(cumLines, numPages, maxLinesPerPage)

	var pages []Page
	for i := 0; i < len(bestSplit)-1; i++ {
		from := bestSplit[i]
		to := bestSplit[i+1]
		if from >= to {
			continue
		}
		pageParas := paras[from:to]
		pages = append(pages, Page{Lines: buildPageLines(pageParas)})
	}

	for i := range pages {
		pages[i].Index = i + 1
		pages[i].Total = len(pages)
	}
	return pages
}

// findOptimalSplit 找到最优的段落分割方案
// cumLines[i] = 前i个段落的累积行数
// 返回分割点数组，如 [0, 5, 10] 表示第1页=段落0-4，第2页=段落5-9
func findOptimalSplit(cumLines []int, numPages, maxLines int) []int {
	n := len(cumLines) - 1 // 段落总数
	target := cumLines[n] / numPages

	// 对于2-3页的情况，直接枚举找最优解
	if numPages == 2 {
		return findBestSplit2(cumLines, n, maxLines, target)
	}

	// 通用情况：贪心 + 偏向后移
	splits := []int{0}
	for p := 0; p < numPages-1; p++ {
		from := splits[len(splits)-1]
		remaining := numPages - p
		remainingLines := cumLines[n] - pageLines(cumLines, from, from)
		thisTarget := remainingLines / remaining

		bestIdx := from + 1
		bestDiff := cumLines[n] // 大数
		for j := from + 1; j <= n; j++ {
			pl := pageLines(cumLines, from, j)
			if pl > maxLines {
				break
			}
			diff := pl - thisTarget
			if diff < 0 {
				diff = -diff
			}
			if diff < bestDiff {
				bestDiff = diff
				bestIdx = j
			}
		}
		splits = append(splits, bestIdx)
	}
	splits = append(splits, n)
	return splits
}

func findBestSplit2(cumLines []int, n, maxLines, target int) []int {
	bestSplit := 1
	bestDiff := cumLines[n]

	for i := 1; i < n; i++ {
		page1 := pageLines(cumLines, 0, i)
		page2 := pageLines(cumLines, i, n)
		if page1 > maxLines || page2 > maxLines {
			continue
		}
		diff := page1 - page2
		if diff < 0 {
			diff = -diff
		}
		if diff < bestDiff {
			bestDiff = diff
			bestSplit = i
		}
	}
	return []int{0, bestSplit, n}
}

// pageLines 计算从段落from到段落to（不含）的显示行数
func pageLines(cumLines []int, from, to int) int {
	if from == 0 {
		return cumLines[to]
	}
	// cumLines包含了段落间空行，需要减去from段落前面的空行
	lines := cumLines[to] - cumLines[from]
	// 如果from > 0，cumLines[from]后面的段落间空行已经算在cumLines[to]里了
	// 但第一个段落前不需要空行，所以要减1
	if from > 0 && to > from {
		lines-- // 去掉页首的段落间空行
	}
	return lines
}

// splitParagraphs 按空行拆分成段落
func splitParagraphs(lines []string) []paragraph {
	var paras []paragraph
	var current []string
	for _, line := range lines {
		if line == "" {
			if len(current) > 0 {
				paras = append(paras, paragraph{lines: current})
				current = nil
			}
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 {
		paras = append(paras, paragraph{lines: current})
	}
	return paras
}

// countDisplayLines 计算段落列表的显示行数（含段落间空行）
func countDisplayLines(paras []paragraph) int {
	if len(paras) == 0 {
		return 0
	}
	total := 0
	for i, p := range paras {
		if i > 0 {
			total++ // 段落间空行
		}
		total += len(p.lines)
	}
	return total
}

// buildPageLines 将段落列表组装成行列表（段落间加空行）
func buildPageLines(paras []paragraph) []string {
	var lines []string
	for i, p := range paras {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, p.lines...)
	}
	return lines
}

func trimEmptyLines(lines []string) []string {
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
