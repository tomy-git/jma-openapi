// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"net/http"

	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/mappers"
)

func (s *Server) GetV1Areas(w http.ResponseWriter, r *http.Request, params gen.GetV1AreasParams) {
	response, err := s.areaUsecase.List(r.Context(), mappers.AreaFilter{
		Parent:     params.Parent,
		Name:       params.Name,
		NameMode:   derefAreaMatchMode(params.NameMatchMode),
		OfficeName: params.OfficeName,
		Child:      params.Child,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func derefAreaMatchMode(mode *gen.AreaMatchMode) gen.AreaMatchMode {
	if mode == nil {
		return gen.AreaMatchMode("exact")
	}

	return *mode
}

func (s *Server) GetV1AreasAreaCode(w http.ResponseWriter, r *http.Request, areaCode string) {
	response, err := s.areaUsecase.Get(r.Context(), areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
