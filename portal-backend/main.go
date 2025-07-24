package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// データモデル
type Dashboard struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Layout    []Gadget  `json:"layout"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Gadget struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position Position               `json:"position"`
	Size     Size                   `json:"size"`
	Config   map[string]interface{} `json:"config"`
	Order    int                    `json:"order"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type GadgetData struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

// インメモリストレージ（本番環境ではデータベース使用を推奨）
type Store struct {
	dashboards map[string]*Dashboard
	mu         sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		dashboards: make(map[string]*Dashboard),
	}
}

func (s *Store) CreateDashboard(dashboard *Dashboard) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dashboard.CreatedAt = time.Now()
	dashboard.UpdatedAt = time.Now()
	s.dashboards[dashboard.ID] = dashboard
	return nil
}

func (s *Store) GetDashboard(id string) (*Dashboard, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dashboard, exists := s.dashboards[id]
	if !exists {
		return nil, fmt.Errorf("dashboard not found")
	}
	return dashboard, nil
}

func (s *Store) UpdateDashboard(dashboard *Dashboard) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.dashboards[dashboard.ID]; !exists {
		return fmt.Errorf("dashboard not found")
	}

	dashboard.UpdatedAt = time.Now()
	s.dashboards[dashboard.ID] = dashboard
	return nil
}

func (s *Store) GetUserDashboards(userID string) ([]*Dashboard, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var dashboards []*Dashboard
	for _, dashboard := range s.dashboards {
		if dashboard.UserID == userID {
			dashboards = append(dashboards, dashboard)
		}
	}

	// 作成日時でソート
	sort.Slice(dashboards, func(i, j int) bool {
		return dashboards[i].CreatedAt.After(dashboards[j].CreatedAt)
	})

	return dashboards, nil
}

// サーバー構造体
type Server struct {
	store *Store
}

func NewServer() *Server {
	return &Server{
		store: NewStore(),
	}
}

// CORS対応ミドルウェア
func (s *Server) enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// JSON レスポンス送信
func (s *Server) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// エラーレスポンス送信
func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ダッシュボード作成
func (s *Server) createDashboard(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dashboard Dashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// IDが指定されていない場合は生成
	if dashboard.ID == "" {
		dashboard.ID = fmt.Sprintf("dash_%d", time.Now().UnixNano())
	}

	if err := s.store.CreateDashboard(&dashboard); err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, dashboard)
}

// ダッシュボード取得
func (s *Server) getDashboard(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "GET" {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスからIDを抽出 /dashboards/{id}
	path := strings.TrimPrefix(r.URL.Path, "/dashboards/")
	if path == "" {
		s.sendError(w, "Dashboard ID required", http.StatusBadRequest)
		return
	}

	dashboard, err := s.store.GetDashboard(path)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusNotFound)
		return
	}

	s.sendJSON(w, dashboard)
}

// ダッシュボード更新
func (s *Server) updateDashboard(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "PUT" {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/dashboards/")
	if path == "" {
		s.sendError(w, "Dashboard ID required", http.StatusBadRequest)
		return
	}

	var dashboard Dashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dashboard.ID = path
	if err := s.store.UpdateDashboard(&dashboard); err != nil {
		s.sendError(w, err.Error(), http.StatusNotFound)
		return
	}

	s.sendJSON(w, dashboard)
}

// ユーザーのダッシュボード一覧取得
func (s *Server) getUserDashboards(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "GET" {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		s.sendError(w, "user_id parameter required", http.StatusBadRequest)
		return
	}

	dashboards, err := s.store.GetUserDashboards(userID)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, dashboards)
}

// ガジェットデータ取得（サンプル実装）
func (s *Server) getGadgetData(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "GET" {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/gadgets/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		s.sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	gadgetType := parts[0]
	gadgetID := parts[1]

	// サンプルデータ生成
	var data interface{}
	switch gadgetType {
	case "weather":
		data = map[string]interface{}{
			"temperature": 22.5,
			"humidity":    65,
			"condition":   "曇り",
			"city":        "東京",
		}
	case "news":
		data = []map[string]interface{}{
			{"title": "サンプルニュース1", "url": "#", "date": "2025-07-24"},
			{"title": "サンプルニュース2", "url": "#", "date": "2025-07-24"},
		}
	case "chart":
		data = map[string]interface{}{
			"labels": []string{"1月", "2月", "3月", "4月", "5月"},
			"values": []int{10, 20, 15, 25, 30},
		}
	default:
		data = map[string]interface{}{
			"message": "Unknown gadget type",
		}
	}

	response := GadgetData{
		ID:   gadgetID,
		Data: data,
	}

	s.sendJSON(w, response)
}

// ルートハンドラー設定
func (s *Server) setupRoutes() {
	http.HandleFunc("/dashboards", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.createDashboard(w, r)
		case "GET":
			s.getUserDashboards(w, r)
		default:
			s.enableCORS(w, r)
			s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.getDashboard(w, r)
		case "PUT":
			s.updateDashboard(w, r)
		default:
			s.enableCORS(w, r)
			s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/gadgets/", s.getGadgetData)
}

// サンプルデータ初期化
func (s *Server) initSampleData() {
	sampleDashboard := &Dashboard{
		ID:     "sample_dashboard",
		UserID: "user1",
		Name:   "サンプルダッシュボード",
		Layout: []Gadget{
			{
				ID:       "weather_1",
				Type:     "weather",
				Position: Position{X: 0, Y: 0},
				Size:     Size{Width: 2, Height: 2},
				Config:   map[string]interface{}{"city": "Tokyo"},
				Order:    0,
			},
			{
				ID:       "news_1",
				Type:     "news",
				Position: Position{X: 2, Y: 0},
				Size:     Size{Width: 4, Height: 3},
				Config:   map[string]interface{}{"category": "tech"},
				Order:    1,
			},
			{
				ID:       "chart_1",
				Type:     "chart",
				Position: Position{X: 0, Y: 2},
				Size:     Size{Width: 3, Height: 2},
				Config:   map[string]interface{}{"type": "line"},
				Order:    2,
			},
		},
	}

	s.store.CreateDashboard(sampleDashboard)
}

func main() {
	server := NewServer()
	server.setupRoutes()
	server.initSampleData()

	port := ":8080"
	log.Printf("サーバーを開始しています ポート%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("サーバー開始エラー:", err)
	}
}
