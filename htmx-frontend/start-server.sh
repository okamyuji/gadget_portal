#!/bin/bash

# HTMX フロントエンド開発サーバー起動スクリプト

PORT=${1:-3000}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 HTMX フロントエンド開発サーバーを起動中..."
echo "📁 ディレクトリ: $SCRIPT_DIR"
echo "🌐 ポート: $PORT"
echo "🔗 URL: http://localhost:$PORT"
echo ""
echo "⚠️  注意: Goバックエンドサーバー (http://localhost:8080) が起動していることを確認してください"
echo ""

# Pythonを使用してシンプルなHTTPサーバーを起動
if command -v python3 &> /dev/null; then
    echo "Python3を使用してサーバーを起動します..."
    cd "$SCRIPT_DIR"
    python3 -m http.server "$PORT"
elif command -v python &> /dev/null; then
    echo "Python2を使用してサーバーを起動します..."
    cd "$SCRIPT_DIR"
    python -m SimpleHTTPServer "$PORT"
elif command -v npx &> /dev/null; then
    echo "npx serve を使用してサーバーを起動します..."
    cd "$SCRIPT_DIR"
    npx serve -l "$PORT"
else
    echo "❌ エラー: Python または npx が見つかりません"
    echo "以下のいずれかをインストールしてください:"
    echo "  - Python 3: https://www.python.org/"
    echo "  - Node.js (npx): https://nodejs.org/"
    exit 1
fi