// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"net/http"
	"strings"

	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/mappers"
)

func (s *Server) GetV1Areas(w http.ResponseWriter, r *http.Request, params gen.GetV1AreasParams) {
	response, err := s.areaUsecase.List(r.Context(), mappers.AreaFilter{
		Parent:     emptyStringToNil(params.Parent),
		Name:       emptyStringToNil(params.Name),
		NameMode:   derefAreaMatchMode(params.NameMatchMode),
		OfficeName: emptyStringToNil(params.OfficeName),
		Child:      emptyStringToNil(params.Child),
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

func emptyStringToNil(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func (s *Server) GetV1AreasAreaCode(w http.ResponseWriter, r *http.Request, areaCode string) {
	response, err := s.areaUsecase.Get(r.Context(), areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
