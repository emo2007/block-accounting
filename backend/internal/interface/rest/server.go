package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	*chi.Mux

	ctx context.Context

	log         *slog.Logger
	addr        string
	tls         bool
	controllers *controllers.RootController

	closeMu sync.RWMutex
	closed  bool
}

func NewServer(
	log *slog.Logger,
	conf config.RestConfig,
	controllers *controllers.RootController,
) *Server {
	s := &Server{
		log:         log,
		addr:        conf.Address,
		tls:         conf.TLS,
		controllers: controllers,
	}

	s.buildRouter()

	return s
}

func (s *Server) Serve(ctx context.Context) error {
	s.ctx = ctx

	s.log.Info(
		"starting rest interface",
		slog.String("addr", s.addr),
		slog.Bool("tls", s.tls),
	)

	if s.tls {
		return http.ListenAndServeTLS(s.addr, "/todo", "/todo", s)
	}

	return http.ListenAndServe(s.addr, s)
}

func (s *Server) Close() {
	s.closeMu.Lock()

	s.closed = true

	s.closeMu.Unlock()
}

func (s *Server) buildRouter() {
	s.Mux = chi.NewRouter()

	s.Use(mw.Recoverer)
	s.Use(mw.RequestID)
	s.Use(s.handleMw)
	s.Use(render.SetContentType(render.ContentTypeJSON))

	s.Get("/ping", s.handlePing) // debug

	// auth
	s.Post("/join", s.handleJoin) // new user
	s.Post("/login", nil)         // login

	s.Route("/organization/{organization_id}", func(r chi.Router) {
		s.Route("/transactions", func(r chi.Router) {
			r.Get("/", nil)           // list
			r.Post("/", nil)          // add
			r.Put("/{tx_id}", nil)    // update / approve (or maybe body?)
			r.Delete("/{tx_id}", nil) // remove
		})

		s.Post("/invite", nil) // create a new invite link

		s.Route("/employees", func(r chi.Router) {
			r.Get("/", nil)                 // list
			r.Post("/", nil)                // add
			r.Put("/{employee_id}", nil)    // update (or maybe body?)
			r.Delete("/{employee_id}", nil) // remove
		})
	})

}

func (s *Server) responseError(w http.ResponseWriter, e error) {
	apiErr := mapError(e)

	out, err := json.Marshal(apiErr)
	if err != nil {
		s.log.Error("error marshal api error", logger.Err(err))

		return
	}

	w.WriteHeader(apiErr.Code)
	w.Write(out)
}

func (s *Server) handleMw(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s.closeMu.RLock()
		defer s.closeMu.RUnlock()

		if s.closed { // keep mutex closed
			return
		}

		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) handleJoin(w http.ResponseWriter, req *http.Request) {
	if err := s.controllers.Auth.Join(w, req); err != nil {
		s.responseError(w, err)
	}
}

func (s *Server) handlePing(w http.ResponseWriter, req *http.Request) {
	s.log.Debug("ping request")

	if err := s.controllers.Ping.Ping(w, req); err != nil {
		s.responseError(w, err)
	}
}
