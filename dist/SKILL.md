---
name: text2img
description: Convert text content into beautifully formatted multi-page screenshot images, with smart pagination and even text distribution. Use this skill when users want to turn text into images for social media (Xiaohongshu/小红书, WeChat Moments, etc.), create text screenshots, or generate paginated text images. Triggers on requests like "convert this text to images", "make screenshots of this text", "generate text images for Xiaohongshu", or any request involving turning text content into picture format.
---

# text2img

CLI tool at `text2img` (installed in PATH). Converts text to evenly-paginated screenshot images.

## Workflow

1. Write user's text content to `/tmp/text2img_input.txt`
2. Run `text2img` with appropriate options
3. Report generated image paths to user

## Command

```bash
text2img -i <input-file> -o <output-dir> [-t template] [-size fontsize] [-lh lineheight] [-font path]
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-i` | stdin | Input text file path |
| `-o` | `./output` | Output directory |
| `-t` | `default` | Template: `default` (white), `dark`, `warm` |
| `-size` | 36 | Font size in pixels |
| `-lh` | 1.8 | Line height multiplier |
| `-font` | system | Custom font file path |
| `-list` | - | List available templates |

## Examples

```bash
# Default white template
text2img -i /tmp/text2img_input.txt -o ./text2img_output

# Dark template
text2img -i /tmp/text2img_input.txt -t dark -o ./text2img_output

# Custom font size
text2img -i /tmp/text2img_input.txt -size 42 -o ./text2img_output

# Pipe input
echo "Hello World" | text2img -t warm -o ./text2img_output
```

## Notes

- Output images are 1080x1440 PNG (portrait, optimized for mobile)
- Text is auto-wrapped, paginated at paragraph boundaries, and vertically distributed
- Page numbers shown automatically when content spans multiple pages
- System Chinese fonts are auto-detected (PingFang on macOS, msyh on Windows, Noto on Linux)
