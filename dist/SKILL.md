---
name: text2img
description: 将文字内容转换为精美的分页截图图片，支持智能分页和文字均匀分布。当用户需要将文字转为图片发小红书、朋友圈等社交平台时使用。触发场景包括"把这段文字转成截图"、"生成文字图片"、"帮我做小红书图文"等涉及文字转图片的请求。
---

# text2img

将文字转换为均匀分页截图的命令行工具。

## 第0步：确认工具已安装

首次使用前检查 `text2img` 是否存在，不存在则自动安装：

```bash
which text2img || bash ~/.claude/skills/text2img/install.sh
```

## 工作流程

1. 将用户提供的文字内容写入 `/tmp/text2img_input.txt`
2. 调用 `text2img` 生成图片
3. 将生成的图片路径返回给用户

## 命令格式

```bash
text2img -i <输入文件> -o <输出目录> [-t 模板] [-size 字号] [-lh 行高] [-font 字体路径]
```

## 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-i` | stdin | 输入文本文件路径 |
| `-o` | `./output` | 输出目录 |
| `-t` | `default` | 模板：`default`（白底）、`dark`（暗色）、`warm`（暖色） |
| `-size` | 36 | 字体大小（像素） |
| `-lh` | 1.8 | 行高倍数 |
| `-font` | 系统字体 | 自定义字体文件路径 |
| `-list` | - | 列出所有可用模板 |

## 示例

```bash
text2img -i /tmp/text2img_input.txt -o ./text2img_output
text2img -i /tmp/text2img_input.txt -t dark -o ./text2img_output
echo "你好世界" | text2img -t warm -o ./text2img_output
```

## 说明

- 输出：1080x1440 PNG 竖屏图片，适配手机屏幕
- 自动换行、按段落边界智能分页、文字垂直均匀分布
- 多页时自动显示页码
- 自动检测系统中文字体（macOS 苹方 / Windows 微软雅黑 / Linux Noto）
