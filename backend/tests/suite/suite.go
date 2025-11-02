package suite

import (
	"context"
	"log/slog"
	"multibank/backend/internal/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	httpserver "multibank/backend/internal/http-server"

	"multibank/backend/internal/auth"
	"multibank/backend/internal/auth/jwt"
	"multibank/backend/internal/config"
	usersvc "multibank/backend/internal/service/user"
	"multibank/backend/internal/storage/sqlite"
)

type Suite struct {
	*testing.T
	Ctx     context.Context
	Cancel  context.CancelFunc
	Log     *slog.Logger
	Cfg     *config.Config
	JWT     *jwt.Manager
	User    *usersvc.Service
	Auth    *auth.Service
	Server  *httptest.Server
	BaseURL string
	Client  *http.Client
	Storage *sqlite.Storage
}

func New(t *testing.T) *Suite {
	t.Helper()

	cfg, repoRoot := mustLoadTestConfig(t)

	// Преобразуем StoragePath в абсолютный путь относительно корня репо,
	// если в конфиге он относительный (./storage/multibank.db)
	if !filepath.IsAbs(cfg.StoragePath) {
		cfg.StoragePath = filepath.Join(repoRoot, cfg.StoragePath)
	}
	// Создаём каталог БД, если его нет
	if err := os.MkdirAll(filepath.Dir(cfg.StoragePath), 0o755); err != nil {
		t.Fatalf("mkdir for db: %v", err)
	}

	// Storage: используем БД из конфигурации (теперь уже абсолютный путь)
	st, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		t.Fatalf("init sqlite: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)

	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	log := logger.Setup(cfg.Level)

	log = log.With(
		slog.String("scope", "test"),
	)

	log.Info("suite: using config",
		slog.String("env", cfg.Env),
		slog.String("db", cfg.StoragePath),
	)

	repo := sqlite.NewUserRepo(st.DB())
	user := usersvc.New(log, repo)
	j := jwt.New(cfg.HTTPServer.JWTSecret, cfg.HTTPServer.TokenTTL)
	a := auth.New(log, user, j)

	srv := httpserver.New(httpserver.Deps{
		Logger:      log,
		UserService: user,
		AuthService: a,
		JWT:         j,
	}, httpserver.Options{
		RequestTimeout: cfg.HTTPServer.Timeout,
	})
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	return &Suite{
		T:       t,
		Ctx:     ctx,
		Cancel:  cancel,
		Log:     log,
		Cfg:     cfg,
		JWT:     j,
		User:    user,
		Auth:    a,
		Server:  ts,
		BaseURL: ts.URL,
		Client:  ts.Client(),
		Storage: st,
	}
}

// грузим config/local.yaml относительно корня репозитория
func mustLoadTestConfig(t *testing.T) (*config.Config, string) {
	t.Helper()
	// Абсолютный путь к текущему файлу (tests/suite/suite.go)
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// Корень репо = <папка с этим файлом>/../..
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	cfgPath := filepath.Join(repoRoot, "config", "local.yaml")

	cfg := config.MustLoadByPath(cfgPath)
	return cfg, repoRoot
}
