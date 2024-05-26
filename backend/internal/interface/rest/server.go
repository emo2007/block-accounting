package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/pkg/metrics"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	jwt jwt.JWTInteractor,
) *Server {
	s := &Server{
		log:         log,
		addr:        conf.Address,
		tls:         conf.TLS,
		controllers: controllers,
		jwt:         jwt,
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
	router.Get("/refresh", s.handle(s.controllers.Auth.Refresh, "refresh"))

	// open invite link
	router.Get("/invite/{hash}", s.handle(s.controllers.Auth.InviteGet, "invite_open"))
	// join via invite link
	router.Post("/invite/{hash}/join", s.handle(s.controllers.Auth.JoinWithInvite, "invite_join"))

	router.Route("/organizations", func(r chi.Router) {
		r = r.With(s.withAuthorization)

		r.Get("/", s.handle(s.controllers.Organizations.ListOrganizations, "list_organizations"))
		r.Post("/", s.handle(s.controllers.Organizations.NewOrganization, "new_organization"))

		r.Route("/{organization_id}", func(r chi.Router) {
			// Deprecated??
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", s.handle(s.controllers.Transactions.List, "tx_list"))
				r.Post("/", s.handle(s.controllers.Transactions.New, "new_tx"))
				r.Put(
					"/{tx_id}",
					s.handle(s.controllers.Transactions.UpdateStatus, "update_tx_status"),
				)
			})

			r.Route("/payrolls", func(r chi.Router) {
				r.Get("/", s.handle(s.controllers.Transactions.ListPayrolls, "list_payrolls"))
				r.Post("/", s.handle(s.controllers.Transactions.NewPayroll, "new_payroll"))
			})

			r.Route("/multisig", func(r chi.Router) {
				r.Post("/", s.handle(s.controllers.Transactions.NewMultisig, "new_multisig"))
				r.Get("/", s.handle(s.controllers.Transactions.ListMultisigs, "list_multisig"))
			})

			r.Route("/license", func(r chi.Router) {
				r.Get("/", nil)  // list license
				r.Post("/", nil) // deploy contract
			})

			r.Route("/participants", func(r chi.Router) {
				r.Get("/", s.handle(s.controllers.Participants.List, "participants_list"))
				r.Post("/", s.handle(s.controllers.Participants.New, "new_participant"))

				// generate new invite link
				r.Post("/invite", s.handle(s.controllers.Auth.Invite, "invite"))

				r.Route("/{participant_id}", func(r chi.Router) {
					r.Get("/", nil) // todo если успею
				})
			})
		})
	})

	s.Mux = router
}

func (s *Server) handle(
	h func(w http.ResponseWriter, req *http.Request) ([]byte, error),
	method_name string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// started := time.Now()
		// defer func() {
		// 	metrics.RequestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
		// 		time.Since(started).Seconds(), prometheus.Labels{
		// 			"reqId":       fmt.Sprint(r.Context().Value(mw.RequestIDKey)),
		// 			"method_name": method_name,
		// 		},
		// 	)
		// }()

		out, err := h(w, r)
		if err != nil {
			s.log.Error(
				"http error",
				slog.String("method_name", method_name),
				logger.Err(err),
			)

			s.responseError(w, err)

			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(out); err != nil {
			s.log.Error(
				"error write http response",
				slog.String("method_name", method_name),
				logger.Err(err),
			)
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
	// todo add rate limiter && cirquit braker

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

		ctx := ctxmeta.UserContext(r.Context(), user)

		if organizationID := chi.URLParam(r, "organization_id"); organizationID != "" {
			organizationUUID, err := uuid.Parse(organizationID)
			if err != nil {
				s.log.Warn(
					"invalid path org id",
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("endpoint", r.RequestURI),
					slog.String("org path id", organizationID),
					logger.Err(err),
				)

				s.responseError(w, ErrorBadPathParams)
				return
			}

			ctx = ctxmeta.OrganizationIdContext(ctx, organizationUUID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
