package route

import (
	"eMobile/internal/config"
	v1 "eMobile/internal/route/api/v1"
	"eMobile/internal/service"
	"eMobile/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	s    service.Service
	l    logging.Logger
	conf *config.Config
}

type Deps struct {
	Service service.Service
	Logger  logging.Logger
	Config  *config.Config
}

func NewHandler(d Deps) *Handler {
	return &Handler{
		s:    d.Service,
		l:    d.Logger,
		conf: d.Config,
	}
}

func (h *Handler) Init(r *httprouter.Router) {
	hv1 := v1.NewHandler(v1.Deps{
		Service: h.s,
		Logger:  h.l,
		Config:  h.conf,
	})
	hv1.Init(r)
}
