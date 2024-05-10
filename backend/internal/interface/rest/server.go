package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/pkg/metrics"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	*chi.Mux

	ctx context.Context

	log         *slog.Logger
	addr        string
	tls         bool
	controllers *controllers.RootController

	jwt jwt.JWTInteractor

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

	metrics.Initialize(s.Mux)

	return http.ListenAndServe(s.addr, s)
}

func (s *Server) Close() {
	s.closeMu.Lock()

	s.closed = true

	s.closeMu.Unlock()
}

func (s *Server) buildRouter() {
	router := chi.NewRouter()

	router.Use(mw.Recoverer)
	router.Use(mw.RequestID)
	router.Use(s.handleMw)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/ping", s.handle(s.controllers.Ping.Ping, "ping"))

	router.Post("/join", s.handle(s.controllers.Auth.Join, "join"))
	router.Post("/login", s.handle(s.controllers.Auth.Login, "login"))

	router.Route("/organization", func(r chi.Router) {
		r.With(s.withAuthorization)

		r.Get("/", s.handle(s.controllers.Auth.Invite, "organization"))

		r.Route("/{organization_id}", func(r chi.Router) {
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", nil)           // list
				r.Post("/", nil)          // add
				r.Put("/{tx_id}", nil)    // update / approve (or maybe body?)
				r.Delete("/{tx_id}", nil) // remove
			})

			r.Post("/invite/{hash}", s.handle(s.controllers.Auth.Invite, "invite")) // create a new invite link

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", nil)                 // list
				r.Post("/", nil)                // add
				r.Put("/{employee_id}", nil)    // update (or maybe body?)
				r.Delete("/{employee_id}", nil) // remove
			})
		})
	})

	s.Mux = router
}

func (s *Server) handle(
	h func(w http.ResponseWriter, req *http.Request) error,
	method_name string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		defer func() {
			reqId := r.Context().Value(mw.RequestIDKey)

			metrics.RequestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
				time.Since(started).Seconds(), prometheus.Labels{
					"reqId":       fmt.Sprint(reqId),
					"method_name": method_name,
				},
			)
		}()

		if err := h(w, r); err != nil {
			s.log.Error(
				"http error",
				slog.String("method_name", method_name),
				logger.Err(err),
			)

			s.responseError(w, err)
		}
	}
}

func (s *Server) responseError(w http.ResponseWriter, e error) {
	s.log.Error("error handle request", logger.Err(e))

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

func (s *Server) withAuthorization(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenStringRaw := r.Header.Get("Authorization")
		if tokenStringRaw == "" {
			s.log.Warn(
				"unauthorized request",
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("endpoint", r.RequestURI),
			)

			w.WriteHeader(401)

			return
		}

		tokenString := strings.Split(tokenStringRaw, " ")[1]

		user, err := s.jwt.User(tokenString)
		if err != nil {
			s.log.Warn(
				"unauthorized request",
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("endpoint", r.RequestURI),
				logger.Err(err),
			)

			s.responseError(w, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(
			ctxmeta.UserContext(r.Context(), user),
		))
	}

	return http.HandlerFunc(fn)
}
