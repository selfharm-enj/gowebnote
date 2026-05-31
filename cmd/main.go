package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webnote/internal/adapters/postgres"
	auth "webnote/internal/auth/http"
	htmlrender "webnote/internal/htmlpages"
	uihandler "webnote/internal/htmlpages/transport/http"
	noterepo "webnote/internal/note/repository"
	notehttp "webnote/internal/note/transport/http"
	noteuc "webnote/internal/note/usecase"
	userrepo "webnote/internal/user/repository"
	userhttp "webnote/internal/user/transport/http"
	useruc "webnote/internal/user/usecase"
	"webnote/utils"
)

type Api struct {
}

func (a *Api) mount() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	// Postgres init
	pgConf := postgres.MustNewConfig()
	pool := postgres.MustNewPool(context.Background(), pgConf)

	slog.Debug("pool created")

	// Env vars.
	srvPort := os.Getenv("PORT")
	httpType := os.Getenv("HTTP_TYPE")
	if srvPort == "" {
		srvPort = "8080"
	}
	if httpType == "" {
		httpType = "http"
	}

	// Server MUX
	mux := http.NewServeMux()

	// User service
	userRepo := userrepo.NewPostgresUserRepository(pool)
	userHasher := utils.NewBcryptHasher()
	userService := useruc.NewUserService(userRepo, userHasher)
	userHandler := userhttp.NewUserHandler(userService)

	// User Routes
	mux.HandleFunc("POST /api/v1/register/{$}", userHandler.CreateUser)
	mux.HandleFunc("POST /api/v1/login/{$}", userHandler.Login)
	mux.HandleFunc("POST /api/v1/logout/{$}", userHandler.Logout)

	// Note service
	noteRepo := noterepo.NewPostgresNoteRepository(pool)
	noteService := noteuc.NewNoteService(noteRepo)
	noteHandler := notehttp.NewNoteHandler(noteService)

	// Note Routes
	mux.Handle("GET /api/v1/notes/{$}", auth.RequireAuth(noteHandler.GetMyNotes, userService))
	mux.Handle("GET /api/v1/notes/{id}/{$}", auth.RequireAuth(noteHandler.GetNote, userService))
	mux.Handle("PATCH /api/v1/notes/{$}", auth.RequireAuth(noteHandler.UpdateNoteByID, userService))
	mux.Handle("POST /api/v1/notes/{$}", auth.RequireAuth(noteHandler.CreateNote, userService))
	mux.Handle("DELETE /api/v1/notes/{id}/{$}", auth.RequireAuth(noteHandler.DeleteNoteByID, userService))

	// WebUI init
	ui, err := NewUI(srvPort, httpType)
	if err != nil {
		panic(err)
	}
	mux.HandleFunc("GET /{$}", ui.Handler.HomePage)
	mux.HandleFunc("GET /register/{$}", ui.Handler.RegisterPage)
	mux.HandleFunc("GET /login/{$}", ui.Handler.LoginPage)
	mux.HandleFunc("GET /notes/{$}", auth.RequireAuth(ui.Handler.NotesPage, userService))
	mux.HandleFunc("GET /notes/{id}/{$}", auth.RequireAuth(ui.Handler.NotePage, userService))
	mux.HandleFunc("GET /notes/create/{$}", ui.Handler.CreateNotePage)

	// Health check
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "alive"}) // #nosec G104
	})

	// Server Run
	slog.Debug("app started")
	server := &http.Server{
		Addr:              ":" + srvPort,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}
	serverErr := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case <-rootCtx.Done():
		slog.Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			slog.Info("http server failed", "err", err)
		}
		return
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Info("http shutdown with error", "err", err)
		server.Close() // #nosec G104
	}

	pool.Close()
	slog.Info("graceful shutdown was completed")
}

type UI struct {
	Renderer *htmlrender.Renderer
	Handler  *uihandler.UIHandler
}

func NewUI(srvPort, httpType string) (*UI, error) {
	r, err := htmlrender.NewRenderer()
	if err != nil {
		return nil, err
	}
	h := uihandler.NewUIHandler(r, srvPort, httpType)
	return &UI{Renderer: r, Handler: h}, nil
}

func main() {
	api := Api{}
	api.mount()
}
