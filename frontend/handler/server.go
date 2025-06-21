package handler

import (
	"html/template"
	"net/http"

	"github.com/LamThanhNguyen/yoyo-store-backend/frontend/util"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	config        util.RuntimeConfig
	templateCache map[string]*template.Template
	DB            models.DBModel
	Session       *scs.SessionManager
	router        http.Handler
}

func NewServer(
	config util.RuntimeConfig,
	templateCache map[string]*template.Template,
	db models.DBModel,
	session *scs.SessionManager,
) (*Server, error) {
	return &Server{
		config:        config,
		templateCache: templateCache,
		DB:            db,
		Session:       session,
	}, nil
}

func (server *Server) SetupRouter() {
	mux := chi.NewRouter()
	mux.Use(server.SessionLoad)

	mux.Get("/", server.Home)
	mux.Get("/ws", server.WsEndPoint)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(server.Auth)
		mux.Get("/virtual-terminal", server.VirtualTerminal)
		mux.Get("/all-sales", server.AllSales)
		mux.Get("/all-subscriptions", server.AllSubscriptions)
		mux.Get("/sales/{id}", server.ShowSale)
		mux.Get("/subscriptions/{id}", server.ShowSubscription)
		mux.Get("/all-users", server.AllUsers)
		mux.Get("/all-users/{id}", server.OneUser)
	})

	mux.Get("/yoyo/{id}", server.ChargeOnce)
	mux.Post("/payment-succeeded", server.PaymentSucceeded)
	mux.Get("/receipt", server.Receipt)

	mux.Get("/plans/bronze", server.BronzePlan)
	mux.Get("/receipt/bronze", server.BronzePlanReceipt)

	// auth routes
	mux.Get("/login", server.LoginPage)
	mux.Post("/login", server.PostLoginPage)
	mux.Get("/logout", server.Logout)
	mux.Get("/forgot-password", server.ForgotPassword)
	mux.Get("/reset-password", server.ShowResetPassword)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	server.router = mux
}

func (server *Server) Router() http.Handler {
	return server.router
}
