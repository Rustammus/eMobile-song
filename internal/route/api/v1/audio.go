package v1

import (
	"database/sql"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/schema"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) initAudioHandler(r *httprouter.Router) {
	r.POST("/api/v1/audios", h.audioCreate)
	r.GET("/api/v1/audios", h.audioList)
	r.GET("/api/v1/audios/:uuid", h.audioFindByUUID)
	r.PATCH("/api/v1/audios/:uuid", h.audioUpdateByUUID)
	r.DELETE("/api/v1/audios/:uuid", h.audioDeleteByUUID)

	r.GET("/api/v1/audios/:uuid/lyrics", h.audioLyricsList)

}

// audioCreate godoc
// @Tags         Audio API
// @Summary      Create audio
// @Description  Create audio
// @Accept       json
// @Produce      json
// @Param Audio body schema.RequestAudioCreate false "Audio base"
// @Success      201  {object}  ResponseBase[schema.ResponseUUID]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios [post]
func (h *Handler) audioCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	audio := schema.RequestAudioCreate{}

	err := json.NewDecoder(r.Body).Decode(&audio)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read body err")
		return
	}

	audioDTO, err := audio.ToDTO()
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "validation err")
		return
	}

	uuid, err := h.s.Audio.Create(audioDTO)
	if err != nil {
		WriteResponseErr(w, http.StatusInternalServerError, err, "create audio err")
		return
	}

	WriteResponse(w, http.StatusCreated, schema.ResponseUUID{uuid}, "audio created correctly")
}

// audioList godoc
// @Tags         Audio API
// @Summary      List audio by Filter
// @Description  List audio by Filter
// @Accept       json
// @Produce      json
// @Param group 	query string false "exact search"
// @Param song 		query string false "full-text-search (english)"
// @Param after 	query string false "after(include) search"
// @Param before 	query string false "before(include) search"
// @Param link 		query string false "exact search"
// @Param lyric 	query string false "full-text-search (english)"
// @Param limit 	query int false "rows limit"
// @Param offset 	query int false "rows offset"
// @Success      200  {object}  ResponseBasePaginated[schema.ResponseAudioRead]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios [get]
func (h *Handler) audioList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filter := schema.RequestAudioFilter{}
	filter.ScanQuery(r.URL)
	filterDTO, err := filter.ToDTO()
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "validation err")
		return
	}

	pag := h.getPagination(r.URL)

	var audios []dto.AudioRead
	if filterDTO == nil {
		audios, err = h.s.Audio.ListPag(pag)
	} else {
		audios, err = h.s.Audio.ListByFilter(filterDTO, pag)
	}

	nextPag := crud.Pagination{
		Offset: pag.Offset + pag.Limit,
		Limit:  pag.Limit,
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteResponsePaginated(w, http.StatusOK, nextPag, []dto.AudioRead{}, "no rows find")
			return
		}
		WriteResponseErr(w, http.StatusInternalServerError, err, "error on list audios")
		return
	}

	audioSchemas := make([]schema.ResponseAudioRead, 0, len(audios))
	for i := 0; i < len(audios); i++ {
		a := schema.ResponseAudioRead{}
		a.FromDTO(&audios[i])
		audioSchemas = append(audioSchemas, a)
	}

	WriteResponsePaginated(w, http.StatusOK, nextPag, audioSchemas, "audios got correctly")
}

// audioFindByUUID godoc
// @Tags         Audio API
// @Summary      Find by UUID
// @Description  Find by UUID
// @Accept       json
// @Produce      json
// @Param uuid path string false "Audio UUID"
// @Param full query boolean false "With lyrics or not"
// @Success      200  {object}  ResponseBase[schema.ResponseAudioReadFull]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios/{uuid} [get]
func (h *Handler) audioFindByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid, err := h.getUUIDParam(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "invalid uuid in path param")
		return
	}

	audioSchema := schema.ResponseAudioReadFull{}
	if r.URL.Query().Get("full") == "true" {
		audioFull, err := h.s.Audio.FindWithLyric(uuid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteResponse(w, http.StatusOK, struct{}{}, "no rows find")
				return
			}
			WriteResponseErr(w, http.StatusInternalServerError, err, "error on find audio by uuid with lyric")
			return
		}
		audioSchema.FromDTOFull(audioFull)
	} else {
		audio, err := h.s.Audio.Find(uuid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteResponse(w, http.StatusOK, struct{}{}, "no rows find")
				return
			}
			WriteResponseErr(w, http.StatusInternalServerError, err, "error on find audio by uuid")
			return
		}
		audioSchema.FromDTO(audio)
	}

	WriteResponse(w, http.StatusOK, audioSchema, "audio got correctly")
}

// audioUpdateByUUID godoc
// @Tags         Audio API
// @Summary      Update audio by UUID
// @Description  Update audio by UUID
// @Accept       json
// @Produce      json
// @Param uuid path string false "Audio UUID"
// @Param Audio body schema.RequestAudioUpdate false "Audio update base"
// @Success      200  {object}  ResponseBase[schema.ResponseAudioRead]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios/{uuid} [patch]
func (h *Handler) audioUpdateByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid, err := h.getUUIDParam(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "invalid uuid in path param")
		return
	}

	audio := schema.RequestAudioUpdate{}
	err = json.NewDecoder(r.Body).Decode(&audio)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read body err")
		return
	}

	audioDTO, err := audio.ToDTO()
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "validation error")
		return
	}

	readAudioDTO, err := h.s.Audio.Update(uuid, audioDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteResponse(w, http.StatusOK, struct{}{}, "no rows updated")
			return
		}
		WriteResponseErr(w, http.StatusInternalServerError, err, "error on update audio")
		return
	}

	readAudioSchema := schema.ResponseAudioRead{}
	readAudioSchema.FromDTO(readAudioDTO)

	WriteResponse(w, http.StatusOK, readAudioSchema, "audio updated correctly")
}

// audioDeleteByUUID godoc
// @Tags         Audio API
// @Summary      Delete audio by UUID
// @Description  Delete audio by UUID
// @Accept       json
// @Produce      json
// @Param uuid path string false "Audio UUID"
// @Success      200  {object}  ResponseBase[schema.ResponseUUID]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios/{uuid} [delete]
func (h *Handler) audioDeleteByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid, err := h.getUUIDParam(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "invalid uuid in path param")
		return
	}

	err = h.s.Audio.Delete(uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteResponse(w, http.StatusOK, struct{}{}, "no rows deleted")
			return
		}
		WriteResponseErr(w, http.StatusInternalServerError, err, "error on delete audio by uuid")
		return
	}

	WriteResponse(w, http.StatusOK, schema.ResponseUUID{uuid}, "audio deleted correctly")
}

// audioLyricsList godoc
// @Tags         Audio API
// @Summary      List audio lyrics by UUID
// @Description  List audio lyrics by UUID
// @Accept       json
// @Produce      json
// @Param uuid path string false "Audio UUID"
// @Param limit 	query string false "rows limit"
// @Param offset 	query string false "rows offset"
// @Success      200  {object}  ResponseBasePaginated[schema.ResponseLyricRead]
// @Failure      400  {object}  ResponseBaseErr
// @Failure      500  {object}	ResponseBaseErr
// @Router       /audios/{uuid}/lyrics [get]
func (h *Handler) audioLyricsList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid, err := h.getUUIDParam(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "invalid uuid in path param")
		return
	}

	pag := h.getPagination(r.URL)
	nextPag := crud.Pagination{
		Offset: pag.Offset + pag.Limit,
		Limit:  pag.Limit,
	}

	lyrics, err := h.s.Lyric.ListByAudioPag(uuid, pag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteResponsePaginated(w, http.StatusOK, nextPag, []dto.LyricRead{}, "no rows find")
			return
		}
		WriteResponseErr(w, http.StatusInternalServerError, err, "error on list audio lyrics")
		return
	}

	lyricSchemas := make([]schema.ResponseLyricRead, 0, len(lyrics))
	for i := 0; i < len(lyrics); i++ {
		l := schema.ResponseLyricRead{}
		l.FromDTO(&lyrics[i])
		lyricSchemas = append(lyricSchemas, l)
	}
	WriteResponsePaginated(w, http.StatusOK, nextPag, lyricSchemas, "lyrics got correctly")
}
