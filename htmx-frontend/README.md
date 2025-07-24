# HTMX + Alpine.js ガジェットポータル

このディレクトリには、HTMX + Alpine.jsを使用したガジェットポータルのフロントエンド実装が含まれています。

## 特徴

- **軽量**: 最小限のJavaScript（~25KB）
- **シンプル**: HTMLアトリビュートベースの宣言的記述
- **ドラッグ&ドロップ**: Sortable.jsによる直感的なガジェット配置
- **レスポンシブ**: Tailwind CSSによるモダンなデザイン
- **リアルタイム**: Alpine.jsによるリアクティブな状態管理

## セットアップ

### 前提条件

1. Goバックエンドサーバーが `http://localhost:8080` で起動していること
2. 以下のいずれかがインストールされていること：
   - Python 3
   - Node.js (npx)

### 起動方法

#### 方法1: スクリプトを使用（推奨）

```bash
./start-server.sh
# または別のポートを指定
./start-server.sh 3001
```

#### 方法2: Node.jsを使用

```bash
# 依存関係をインストール
npm install

# 開発サーバーを起動
npm start
```

#### 方法3: Pythonを使用

```bash
# Python 3
python3 -m http.server 3000

# Python 2
python -m SimpleHTTPServer 3000
```

### アクセス

ブラウザで http://localhost:3000 にアクセスしてください。

## ディレクトリ構造

```text
htmx-frontend/
├── index.html          # メインアプリケーション
├── package.json        # Node.js設定
├── start-server.sh     # 起動スクリプト
└── README.md          # このファイル
```

## 使用方法

### ガジェットの追加

1. 右上の「ガジェット追加」ボタンをクリック
2. 追加したいガジェットタイプを選択：
   - **天気**: 気温、湿度、天候情報
   - **ニュース**: 最新ニュース一覧
   - **チャート**: データ可視化

### ガジェットの操作

- **移動**: ドラッグハンドル（矢印アイコン）をドラッグして配置変更
- **設定**: 設定アイコンをクリックして設定変更
- **削除**: ×アイコンをクリックして削除

### 自動保存

ガジェットの配置変更や追加・削除は自動的にサーバーに保存されます。

## 技術スタック

### フロントエンド

- **HTMX 1.9.10**: HTMLベースのAJAX通信
- **Alpine.js 3.x**: 軽量なJavaScriptフレームワーク
- **Sortable.js**: ドラッグ&ドロップライブラリ
- **Tailwind CSS**: ユーティリティファーストCSS

### APIエンドポイント

以下のGoバックエンドAPIを使用します：

- `GET /dashboards/{id}` - ダッシュボード取得
- `PUT /dashboards/{id}` - ダッシュボード更新
- `GET /gadgets/{type}/{id}` - ガジェットデータ取得

## 開発情報

### ファイル構成

- **HTML構造**: セマンティックなマークアップ
- **CSS**: Tailwind CSS + カスタムスタイル
- **JavaScript**: Alpine.js + Sortable.js integration

### 主要なAlpine.jsコンポーネント

```javascript
dashboardApp() {
    return {
        dashboard: { layout: [] },     // ダッシュボード状態
        gadgetData: {},               // ガジェットデータ
        loading: true,                // ローディング状態
        error: null,                  // エラー状態
        // ... メソッド
    }
}
```

### カスタマイズ

#### 新しいガジェットタイプの追加

1. HTMLテンプレートにガジェット表示ロジックを追加
2. `addGadget()` メソッドに新しいタイプを追加
3. バックエンドでガジェットデータエンドポイントを実装

#### スタイルのカスタマイズ

Tailwind CSSクラスまたはカスタムCSSで見た目を変更できます。

## トラブルシューティング

### よくある問題

1. **ガジェットが表示されない**
   - Goバックエンドサーバーが起動しているか確認
   - ブラウザの開発者ツールでエラーを確認

2. **ドラッグ&ドロップが動作しない**
   - JavaScriptエラーがないか確認
   - Sortable.jsが正しく読み込まれているか確認

3. **CORSエラー**
   - バックエンドのCORS設定を確認
   - 同じオリジンからアクセスしているか確認

### デバッグ

ブラウザの開発者ツールのコンソールでデバッグ情報を確認できます：

```javascript
// Alpine.jsアプリケーション状態の確認
$el._x_dataStack[0]

// API通信の確認
// Network タブでHTTPリクエストを監視
```

## 本番環境

本番環境では、適切なWebサーバー（Nginx、Apache等）を使用してください。

例（Nginx設定）:

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /path/to/htmx-frontend;
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://backend:8080/;
    }
}
```

## ライセンス

MIT License
