# Go ガジェットポータル API

このディレクトリには、Goで実装されたガジェットポータルのバックエンドAPIが含まれています。

## 特徴

- **軽量**: 標準ライブラリのみを使用した軽量な実装
- **シンプル**: RESTful API設計
- **リアルタイム**: インメモリストレージによる高速レスポンス
- **CORS対応**: フロントエンドとの連携をサポート
- **スケーラブル**: 構造化されたデータモデル

## セットアップ

### 前提条件

- Go 1.24.5 以上がインストールされていること

### 起動方法

```bash
# ディレクトリに移動
cd portal-backend

# 依存関係の確認
go mod tidy

# サーバーを起動
go run main.go
```

### アクセス

サーバーは `http://localhost:8080` で起動します。

## ディレクトリ構造

```text
portal-backend/
├── main.go            # メインアプリケーション
├── go.mod            # Go モジュール設定
└── README.md         # このファイル
```

## API エンドポイント

### ダッシュボード管理

- `POST /dashboards` - ダッシュボード作成
- `GET /dashboards?user_id={id}` - ユーザーのダッシュボード一覧取得
- `GET /dashboards/{id}` - 特定のダッシュボード取得
- `PUT /dashboards/{id}` - ダッシュボード更新

### ガジェットデータ

- `GET /gadgets/{type}/{id}` - ガジェットデータ取得

### サポートされるガジェットタイプ

- **weather**: 天気情報（気温、湿度、天候、都市）
- **news**: ニュース一覧（タイトル、URL、日付）
- **chart**: チャートデータ（ラベル、数値）

## データモデル

### Dashboard

```go
type Dashboard struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Name      string    `json:"name"`
    Layout    []Gadget  `json:"layout"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Gadget

```go
type Gadget struct {
    ID       string                 `json:"id"`
    Type     string                 `json:"type"`
    Position Position               `json:"position"`
    Size     Size                   `json:"size"`
    Config   map[string]interface{} `json:"config"`
    Order    int                    `json:"order"`
}
```

### Position & Size

```go
type Position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

type Size struct {
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

## 使用例

### ダッシュボード作成

```bash
curl -X POST http://localhost:8080/dashboards \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user1",
    "name": "マイダッシュボード",
    "layout": []
  }'
```

### ガジェットデータ取得

```bash
# 天気データ
curl http://localhost:8080/gadgets/weather/weather_1

# ニュースデータ
curl http://localhost:8080/gadgets/news/news_1

# チャートデータ
curl http://localhost:8080/gadgets/chart/chart_1
```

## 技術スタック

### バックエンド

- **Go 1.24.5**: プログラミング言語
- **net/http**: HTTP サーバー（標準ライブラリ）
- **encoding/json**: JSON エンコーディング（標準ライブラリ）
- **sync**: 並行制御（標準ライブラリ）

### ストレージ

- **インメモリ**: 開発・テスト用（本番環境ではデータベース推奨）

## 開発情報

### プロジェクト構造

- **Server**: HTTPサーバーとルーティング管理
- **Store**: データストレージの抽象化
- **Models**: データ構造定義

### 主要なコンポーネント

```go
type Server struct {
    store *Store  // データストレージ
}

type Store struct {
    dashboards map[string]*Dashboard  // ダッシュボード保存
    mu         sync.RWMutex          // 並行制御
}
```

### CORS設定

すべてのリクエストに対してCORSヘッダーを設定：

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

### カスタマイズ

#### 新しいガジェットタイプの追加

1. `getGadgetData()` メソッドにケースを追加
2. 対応するデータ構造を定義
3. フロントエンドでガジェット表示を実装

#### データベース統合

本番環境では、`Store` インターフェースを実装してデータベースと連携：

```go
type DatabaseStore struct {
    db *sql.DB
}

func (s *DatabaseStore) CreateDashboard(dashboard *Dashboard) error {
    // データベース操作
}
```

## トラブルシューティング

### よくある問題

1. **ポート8080が使用中**
   - `main.go` の `port` 変数を変更
   - または他のプロセスを停止

2. **CORSエラー**
   - フロントエンドのオリジンが許可されているか確認
   - プリフライトリクエスト（OPTIONS）が適切に処理されているか確認

3. **JSONパースエラー**
   - リクエストボディのJSON形式を確認
   - Content-Typeヘッダーが正しく設定されているか確認

### デバッグ

ログの確認：

```bash
# サーバー起動時のメッセージ
サーバーを開始しています ポート:8080

# エラーログ
go run main.go 2>&1 | tee server.log
```

### APIテスト

```bash
# ヘルスチェック（ダッシュボード一覧）
curl http://localhost:8080/dashboards?user_id=user1

# サンプルダッシュボードの確認
curl http://localhost:8080/dashboards/sample_dashboard
```

## 本番環境

本番環境での推奨設定：

### 環境変数

```bash
export PORT=8080
export DB_URL="postgres://user:pass@localhost/gadget_portal"
export CORS_ORIGINS="https://yourdomain.com"
```

### サーバー設定

```go
// main.go での設定例
port := os.Getenv("PORT")
if port == "" {
    port = ":8080"
}
```

### データベース

PostgreSQL設定例：

```sql
CREATE TABLE dashboards (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    layout JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## パフォーマンス

### 最適化ポイント

- データベース接続プールの使用
- キャッシュ層の導入
- ガジェットデータの非同期取得
- ページネーション実装

### モニタリング

```go
// レスポンス時間の計測
start := time.Now()
// ... 処理 ...
log.Printf("API took %v", time.Since(start))
```

## ライセンス

MIT License
