---
name: text2img
description: Convert text content into beautifully formatted multi-page screenshot images, with smart pagination and even text distribution. Use this skill when users want to turn text into images for social media (Xiaohongshu/小红书, WeChat Moments, etc.), create text screenshots, or generate paginated text images. Triggers on requests like "convert this text to images", "make screenshots of this text", "generate text images for Xiaohongshu", or any request involving turning text content into picture format.
---

# text2img

CLI tool that converts text to evenly-paginated screenshot images.

## Step 0: Ensure binary is available

Before first use, check if `text2img` exists:

```bash
which text2img || bash ~/.claude/skills/text2img/install.sh
```

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
text2img -i /tmp/text2img_input.txt -o ./text2img_output
text2img -i /tmp/text2img_input.txt -t dark -o ./text2img_output
echo "Hello World" | text2img -t warm -o ./text2img_output
```

## Notes

- Output: 1080x1440 PNG, portrait, optimized for mobile
- Auto line-wrap, paragraph-boundary pagination, vertical text distribution
- Page numbers auto-shown when multi-page
- Chinese fonts auto-detected per platform (PingFang/msyh/Noto)
