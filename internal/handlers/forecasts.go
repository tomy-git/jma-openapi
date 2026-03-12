// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"net/http"

	"github.com/tomy-git/jma-openapi/internal/gen"
)

func (s *Server) GetV1ForecastsOfficeCode(w http.ResponseWriter, r *http.Request, officeCode string) {
	response, err := s.forecastUsecase.Get(r.Context(), officeCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeAreasAreaCode(w http.ResponseWriter, r *http.Request, officeCode string, areaCode string) {
	response, err := s.forecastUsecase.GetArea(r.Context(), officeCode, areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeAreas(w http.ResponseWriter, r *http.Request, officeCode string) {
	response, err := s.forecastUsecase.ListAreas(r.Context(), officeCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeWeatherAreas(w http.ResponseWriter, r *http.Request, officeCode string) {
	response, err := s.forecastUsecase.ListWeatherAreas(r.Context(), officeCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeWeatherAreasAreaCode(w http.ResponseWriter, r *http.Request, officeCode string, areaCode string) {
	response, err := s.forecastUsecase.GetWeatherArea(r.Context(), officeCode, areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeTemperatureAreas(w http.ResponseWriter, r *http.Request, officeCode string) {
	response, err := s.forecastUsecase.ListTemperatureAreas(r.Context(), officeCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeTemperatureAreasAreaCode(w http.ResponseWriter, r *http.Request, officeCode string, areaCode string) {
	response, err := s.forecastUsecase.GetTemperatureArea(r.Context(), officeCode, areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeAreasResolve(w http.ResponseWriter, r *http.Request, officeCode string, params gen.GetV1ForecastsOfficeCodeAreasResolveParams) {
	response, err := s.forecastUsecase.ResolveAreas(r.Context(), officeCode, params)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeAreasAreaCodeLatest(w http.ResponseWriter, r *http.Request, officeCode string, areaCode string) {
	response, err := s.forecastUsecase.GetAreaLatest(r.Context(), officeCode, areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetV1ForecastsOfficeCodeAreasAreaCodeTimeseries(w http.ResponseWriter, r *http.Request, officeCode string, areaCode string) {
	response, err := s.forecastUsecase.GetAreaTimeseries(r.Context(), officeCode, areaCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
