package server

import (
	"log"
	"net/http"
	"pocketTeleBot/pkg/database"
	"pocketTeleBot/pkg/pocketAPI"
	"strconv"
)

type AuthorizationServer struct {
	server       *http.Server
	pocketClient *pocketAPI.Client
	tokenDb      database.TokenDB
	redirectURL  string
}

func NewAuthorizationServer(pocketClient *pocketAPI.Client, tokenDb database.TokenDB, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenDb: tokenDb, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":8000",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdQuery := r.URL.Query().Get("chat_id")
	if chatIdQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIdQuery, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenDb.Get(chatID, database.RequestToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	resp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("chatID: %d\nrequestToken: %s\naccesToken: %s", chatID, requestToken, resp.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
