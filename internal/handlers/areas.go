// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"net/http"

	"github.com/tomy-git/jma-openapi/internal/gen"
)

func (s *Server) GetV1Areas(w http.ResponseWriter, r *http.Request, params gen.GetV1AreasParams) {
	response, err := s.areaUsecase.List(r.Context(), params.Parent)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1AreasAreaCode(w http.ResponseWriter, r *http.Request, areaCode string) {
	response, err := s.areaUsecase.Get(r.Context(), areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
