package handlers

import "net/http"

func (s *Server) GetHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.healthUsecase.Execute())
}
