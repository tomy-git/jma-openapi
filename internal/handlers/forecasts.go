package handlers

import "net/http"

func (s *Server) GetV1ForecastsOfficeCode(w http.ResponseWriter, r *http.Request, officeCode string) {
	response, err := s.forecastUsecase.Get(r.Context(), officeCode)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
