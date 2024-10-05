package v1

import (
	"eMobile/internal/config"
	"eMobile/internal/crud"
	"eMobile/internal/service"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"net/url"
	"strconv"
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
	h.initAudioHandler(r)
}

func (h *Handler) getPagination(u *url.URL) crud.Pagination {
	q := u.Query()
	offset := q.Get("offset")
	limit := q.Get("limit")

	pag := crud.Pagination{
		Offset: 0,
		Limit:  h.conf.Server.PagLimit,
	}

	offsetNum, err := strconv.Atoi(offset)
	if offsetNum > 0 && err == nil {
		pag.Offset = offsetNum
	}

	limitNum, err := strconv.Atoi(limit)
	if limitNum > 0 && limitNum < h.conf.Server.PagLimit && err == nil {
		pag.Limit = limitNum
	}

	return pag
}

func (h *Handler) getUUIDParam(params httprouter.Params) (pgtype.UUID, error) {
	strUUID := params.ByName("uuid")
	uuid := pgtype.UUID{}
	err := uuid.Scan(strUUID)
	return uuid, err
}
