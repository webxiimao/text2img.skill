# text2img

将大段文字自动转换为精美的分页截图图片，适用于小红书、朋友圈等社交平台。

支持智能分页（段落边界不截断）、文字均匀分布、多种模板样式。

## 功能特性

- 智能分页：基于段落边界分割，保证段落完整性
- 均匀排版：文字自动垂直均匀分布，铺满整页
- 多种模板：白底（default）、暗色（dark）、暖色纸张（warm）
- 跨平台：支持 macOS / Windows / Linux
- 零依赖：单一可执行文件，不需要安装运行时
- Skill 集成：可作为 Claude Code Skill 使用

## 安装

### macOS（Apple Silicon / M系列芯片）

```bash
# 1. 下载或克隆仓库
git clone https://github.com/webxiimao/text2img.skill.git
cd text2img.skill

# 2. 运行安装脚本
bash skills/text2img/install.sh
```

**macOS 安全提示：** 首次运行时系统可能会提示"无法打开，因为无法验证开发者"。解决方法：

```bash
# 方法一：移除隔离属性（推荐）
sudo xattr -d com.apple.quarantine /usr/local/bin/text2img

# 方法二：通过系统设置允许
# 打开「系统设置 → 隐私与安全性」，找到被阻止的 text2img，点击「仍要打开」
```

### macOS（Intel x86）

安装步骤同上，install.sh 会自动检测芯片架构并选择对应的二进制文件。

### Linux

```bash
git clone https://github.com/webxiimao/text2img.skill.git
cd text2img.skill
bash skills/text2img/install.sh
```

> 需要系统已安装中文字体（如 Noto Sans CJK）。若未安装：
> ```bash
> # Ubuntu / Debian
> sudo apt install fonts-noto-cjk
>
> # CentOS / RHEL
> sudo yum install google-noto-sans-cjk-fonts
> ```

### Windows

1. 下载仓库并解压
2. 将 `skills/text2img/text2img-windows-x86.exe` 复制到任意目录
3. 重命名为 `text2img.exe`
4. 将该目录添加到系统 PATH 环境变量

或直接在 PowerShell 中使用完整路径调用。

### 从源码编译（所有平台）

需要 Go 1.20+：

```bash
git clone https://github.com/webxiimao/text2img.skill.git
cd text2img.skill
GOPROXY=https://goproxy.cn,direct go build -o text2img .
# 将生成的 text2img 移动到 PATH 目录下
```

## 使用方式

```bash
# 从文件生成（默认白底模板）
text2img -i input.txt -o ./output

# 从剪贴板生成（macOS）
pbpaste | text2img -o ./output

# 使用暗色模板
text2img -i input.txt -t dark -o ./output

# 使用暖色纸张模板
text2img -i input.txt -t warm -o ./output

# 自定义字号和行高
text2img -i input.txt -size 42 -lh 2.0 -o ./output

# 指定自定义字体
text2img -i input.txt -font ./MyFont.ttf -o ./output

# 查看所有可用模板
text2img -list
```

## 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-i` | stdin | 输入文本文件路径 |
| `-o` | `./output` | 输出目录 |
| `-t` | `default` | 模板：`default`（白底）/ `dark`（暗色）/ `warm`（暖色） |
| `-size` | 36 | 字体大小（像素） |
| `-lh` | 1.8 | 行高倍数 |
| `-font` | 系统字体 | 自定义字体文件路径 |
| `-list` | - | 列出所有可用模板 |

## 输出说明

- 图片尺寸：1080 × 1440 像素（竖屏，适配手机）
- 格式：PNG
- 多页时自动显示页码（如 `1 / 3`）

## 作为 Claude Code Skill 使用

将 `skills/text2img/` 目录复制到 `~/.claude/skills/`：

```bash
cp -r skills/text2img ~/.claude/skills/
```

之后在 Claude Code 中直接说"帮我把这段文字转成截图"即可自动调用。

## 项目结构

```
text2img.skill/
├── config.go                          # 页面配置和模板定义
├── font.go                            # 字体加载（跨平台）
├── layout.go                          # 文字排版和智能分页算法
├── renderer.go                        # 图片渲染引擎
├── main.go                            # CLI 入口
├── go.mod / go.sum                    # Go 依赖
├── skills/text2img/                   # Skill 分发包
│   ├── SKILL.md                       # Skill 定义文件
│   ├── install.sh                     # 一键安装脚本
│   ├── text2img-mac-apple-silicon     # macOS ARM64
│   ├── text2img-mac-intel-x86         # macOS x86_64
│   ├── text2img-linux-x86             # Linux x86_64
│   └── text2img-windows-x86.exe       # Windows x86_64
└── README.md
```

## License

MIT
