package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
)

type envelope map[string]any

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeData(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, envelope{"data": data})
}

func writeList[T any](w http.ResponseWriter, result repositories.ListResult[T], page int, pageSize int) {
	writeJSON(w, http.StatusOK, envelope{
		"data": result.Items,
		"meta": envelope{"page": page, "page_size": pageSize, "total": result.Total},
	})
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	msg := "internal server error"
	switch {
	case errors.Is(err, domain.ErrValidation):
		status, code, msg = http.StatusBadRequest, "VALIDATION_ERROR", err.Error()
	case errors.Is(err, domain.ErrUnauthenticated):
		status, code, msg = http.StatusUnauthorized, "UNAUTHENTICATED", "authentication required"
	case errors.Is(err, domain.ErrForbidden):
		status, code, msg = http.StatusForbidden, "FORBIDDEN", "forbidden"
	case errors.Is(err, domain.ErrNotFound):
		status, code, msg = http.StatusNotFound, "NOT_FOUND", "resource not found"
	case errors.Is(err, domain.ErrConflict):
		status, code, msg = http.StatusConflict, "CONFLICT", "resource conflict"
	}
	writeJSON(w, status, envelope{"error": envelope{"code": code, "message": msg}})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func page(r *http.Request) repositories.PageFilter {
	q := r.URL.Query()
	p, _ := strconv.Atoi(q.Get("page"))
	ps, _ := strconv.Atoi(q.Get("page_size"))
	f := repositories.PageFilter{Page: p, PageSize: ps}
	limit, _ := f.LimitOffset()
	if f.Page < 1 {
		f.Page = 1
	}
	f.PageSize = limit
	return f
}
