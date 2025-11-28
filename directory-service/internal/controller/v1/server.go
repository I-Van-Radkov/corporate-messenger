package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv *http.Server
	db  *pgxpool.Pool
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, db *pgxpool.Pool) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      nil,
	}

	return &Server{
		srv: srv,
		db:  db,
	}
}
