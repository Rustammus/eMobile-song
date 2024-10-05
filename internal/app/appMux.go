package app

import (
	"eMobile/internal/config"
	"eMobile/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type AppMux struct {
	r    *httprouter.Router
	l    logging.Logger
	conf *config.Config
}

type Deps struct {
	Router *httprouter.Router
	Log    logging.Logger
	Conf   *config.Config
}

func NewAppMux(d *Deps) *AppMux {
	return &AppMux{
		r:    d.Router,
		l:    d.Log,
		conf: d.Conf,
	}
}

func (s *AppMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rcv := recover(); rcv != nil {
			s.PanicHandler(w, r, rcv)
		}
	}()
	r.Body = http.MaxBytesReader(w, r.Body, 10<<10)
	s.LogMiddleware(s.r).ServeHTTP(w, r)
}

func (s *AppMux) PanicHandler(w http.ResponseWriter, r *http.Request, rcv any) {
	w.WriteHeader(http.StatusInternalServerError)
	s.l.Error("recover from panic: ", rcv)
}

func (s *AppMux) LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sTime := time.Now().UTC()
		handler.ServeHTTP(w, r)
		eTime := time.Since(sTime)
		s.l.Debugf("[%s]\t%s\t%s\ttime:%s", r.Method, r.URL.Path, r.RemoteAddr, eTime.String())
	})
}
