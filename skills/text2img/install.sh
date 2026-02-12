#!/bin/bash
set -e

echo "=== text2img 安装脚本 ==="
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 检测平台，选择对应二进制
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS-$ARCH" in
  Darwin-arm64)   BINARY="text2img-mac-apple-silicon" ;;
  Darwin-x86_64)  BINARY="text2img-mac-intel-x86" ;;
  Linux-x86_64)   BINARY="text2img-linux-x86" ;;
  Linux-aarch64)  echo "暂不支持 Linux ARM64"; exit 1 ;;
  *)              echo "不支持的平台: $OS $ARCH"; exit 1 ;;
esac

if [ ! -f "$SCRIPT_DIR/$BINARY" ]; then
  echo "未找到二进制文件: $BINARY"
  exit 1
fi

echo "检测到平台: $OS $ARCH -> $BINARY"

# 安装二进制
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  echo "需要管理员权限安装到 $INSTALL_DIR"
  sudo cp "$SCRIPT_DIR/$BINARY" "$INSTALL_DIR/text2img"
  sudo chmod +x "$INSTALL_DIR/text2img"
else
  cp "$SCRIPT_DIR/$BINARY" "$INSTALL_DIR/text2img"
  chmod +x "$INSTALL_DIR/text2img"
fi
echo "已安装: $INSTALL_DIR/text2img"

# 安装 skill
SKILL_DIR="$HOME/.claude/skills/text2img"
mkdir -p "$SKILL_DIR"
cp "$SCRIPT_DIR/SKILL.md" "$SKILL_DIR/SKILL.md"
echo "已安装 Skill: $SKILL_DIR"

echo ""
echo "=== 安装完成 ==="
echo "  text2img -i input.txt -o ./output"
echo "  echo '你好世界' | text2img -t dark -o ./output"
echo "  text2img -list"
