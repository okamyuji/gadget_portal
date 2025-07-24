# ガジェットポータルサイト

ユーザーがドラッグ&ドロップで自由にカスタマイズできるガジェット型ポータルサイトの実装です。

## プロジェクト構成

```text
gadget_portal/
├── portal-backend/         # Go標準ライブラリのみのバックエンド
└── htmx-frontend/          # HTMX + Alpine.js フロントエンド
```

## アーキテクチャ

```text
Frontend (HTMX + Alpine.js) ←→ REST API ←→ Backend (Go標準ライブラリ)
                                               ↓
                                        JSON データストレージ
```

## 特徴

### バックエンド (Go)

- **軽量**: Go標準ライブラリのみ使用
- **REST API**: シンプルで拡張可能なAPI設計
- **CORS対応**: フロントエンドとの連携
- **インメモリストレージ**: 開発用（本番用DB接続例も含む）

### フロントエンド (HTMX + Alpine.js)

- **軽量**: 最小限のJavaScript（~25KB）
- **ドラッグ&ドロップ**: 直感的なガジェット配置
- **レスポンシブ**: モダンなデザイン
- **リアルタイム**: リアクティブな状態管理

## クイックスタート

### 1. バックエンドの起動

```bash
cd portal-backend
go run main.go
```

サーバーが `http://localhost:8080` で起動します。

### 2. フロントエンドの起動

```bash
cd htmx-frontend
./start-server.sh
```

フロントエンドが `http://localhost:3000` で起動します。

### 3. アクセス

ブラウザで http://localhost:3000 にアクセスしてガジェットポータルを使用できます。

## 機能

### ガジェット管理

- **天気ガジェット**: 気温、湿度、天候情報
- **ニュースガジェット**: 最新ニュース一覧
- **チャートガジェット**: データ可視化

### ユーザー操作

- **ドラッグ&ドロップ**: ガジェットの配置変更
- **追加・削除**: ガジェットの動的な管理
- **設定変更**: ガジェット個別の設定
- **自動保存**: 変更の自動保存

## API仕様

### ダッシュボード管理

- `POST /dashboards` - ダッシュボード作成
- `GET /dashboards?user_id={id}` - ユーザーのダッシュボード一覧
- `GET /dashboards/{id}` - ダッシュボード詳細取得
- `PUT /dashboards/{id}` - ダッシュボード更新

### ガジェットデータ

- `GET /gadgets/{type}/{id}` - ガジェットデータ取得

## 技術スタック

### バックエンド

- **Go**: 標準ライブラリのみ
- **HTTP**: net/httpパッケージ
- **JSON**: encoding/jsonパッケージ

### フロントエンド

- **HTMX 1.9.10**: HTMLベースのAJAX
- **Alpine.js 3.x**: 軽量JavaScriptフレームワーク
- **Sortable.js**: ドラッグ&ドロップ
- **Tailwind CSS**: ユーティリティファーストCSS

## 開発

### 前提条件

- Go 1.19以上
- Python 3 または Node.js（フロントエンド開発サーバー用）

### 開発サーバーの起動

1. **バックエンド**:

   ```bash
   cd portal-backend
   go run main.go
   ```

2. **フロントエンド**:

   ```bash
   cd htmx-frontend
   ./start-server.sh
   ```

### カスタマイズ

#### 新しいガジェットタイプの追加

1. **バックエンド**: `getGadgetData` 関数に新しいタイプを追加
2. **フロントエンド**: HTMLテンプレートに表示ロジックを追加

#### スタイルの変更

Tailwind CSSクラスまたはカスタムCSSでデザインをカスタマイズできます。

## 本番環境デプロイ

### Docker使用例

```dockerfile
# バックエンド
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY portal-backend/ .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Nginx設定例

```nginx
upstream backend {
    server app:8080;
}

server {
    listen 80;
    
    location /api/ {
        proxy_pass http://backend/;
    }
    
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
}
```

## パフォーマンス

### バックエンド

- メモリ使用量: ~10MB
- レスポンス時間: <10ms（ローカル）
- 並行接続: 高い（Goのgoroutine）

### フロントエンド

- JavaScript バンドルサイズ: ~25KB
- 初期読み込み時間: <1秒
- ドラッグ&ドロップ: 60FPS

## セキュリティ

- **CORS設定**: 適切なオリジン制御
- **入力検証**: JSONスキーマ検証
- **エラーハンドリング**: 適切なエラーレスポンス

## 今後の拡張予定

- [ ] ユーザー認証システム
- [ ] データベース統合（PostgreSQL/MySQL）
- [ ] Webソケット対応（リアルタイム更新）
- [ ] ガジェット設定UI改善
- [ ] モバイル対応の強化
- [ ] テーマシステム

## トラブルシューティング

### よくある問題

1. **CORSエラー**: バックエンドが正しく起動しているか確認
2. **ポート競合**: 8080/3000ポートが使用可能か確認
3. **ガジェット表示エラー**: ブラウザの開発者ツールでエラー確認

### ログ確認

```bash
# バックエンドログ
cd portal-backend
go run main.go

# フロントエンドログ
# ブラウザの開発者ツールコンソール
```

## ライセンス

MIT License

## 貢献

プルリクエストやIssueの報告を歓迎します。

## サポート

技術的な質問や問題については、GitHub Issues をご利用ください。
