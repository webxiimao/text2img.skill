---
name: text2img
description: Convert text content into beautifully formatted multi-page screenshot images, with smart pagination and even text distribution. Use this skill when users want to turn text into images for social media (Xiaohongshu/小红书, WeChat Moments, etc.), create text screenshots, or generate paginated text images. Triggers on requests like "convert this text to images", "make screenshots of this text", "generate text images for Xiaohongshu", or any request involving turning text content into picture format.
---

# text2img

CLI tool that converts text to evenly-paginated screenshot images.

## Step 0: Ensure binary is available

Before first use, check if `text2img` exists. If not, auto-install:

```bash
# Check
which text2img

# If not found, clone and build (requires Go)
git clone https://github.com/webxiimao/text2img.skill.git /tmp/text2img-build \
  && cd /tmp/text2img-build \
  && GOPROXY=https://goproxy.cn,direct go build -o /usr/local/bin/text2img . \
  && rm -rf /tmp/text2img-build

# If Go is also not installed, install it first:
# macOS: brew install go
# Linux: sudo apt install golang-go / sudo yum install golang
```

If both `go` and `brew`/`apt` are unavailable, inform the user to manually install from: https://github.com/webxiimao/text2img.skill

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

# Pipe input
echo "Hello World" | text2img -t warm -o ./text2img_output
```

## Notes

- Output images are 1080x1440 PNG (portrait, optimized for mobile)
- Text is auto-wrapped, paginated at paragraph boundaries, and vertically distributed
- Page numbers shown automatically when content spans multiple pages
- System Chinese fonts are auto-detected (PingFang on macOS, msyh on Windows, Noto on Linux)
