package handler

import (
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"
	"net/http"
)

func (h Handlers) CreateContactRequest(w http.ResponseWriter, r *http.Request) {
	var req ContactRequestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	c := domain.ContactRequest{FullName: req.FullName, Email: req.Email, Phone: req.Phone, Subject: req.Subject, Message: req.Message, Status: req.Status}
	if err := h.ContactRequests.Create(r.Context(), &c); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, c)
}
func (h Handlers) GetContactRequest(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	c, err := h.ContactRequests.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, c)
}
func (h Handlers) ListContactRequests(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.ContactRequests.List(r.Context(), repositories.ContactRequestFilter{PageFilter: pf, Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeList(w, res, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateContactRequest(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req ContactRequestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	c := domain.ContactRequest{ID: rid, FullName: req.FullName, Email: req.Email, Phone: req.Phone, Subject: req.Subject, Message: req.Message, Status: req.Status}
	if err := h.ContactRequests.Update(r.Context(), &c); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, c)
}
func (h Handlers) DeleteContactRequest(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.ContactRequests.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) ResolveContactRequest(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.ContactRequests.Resolve(r.Context(), rid, actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) CreateTeamMember(w http.ResponseWriter, r *http.Request) {
	var req TeamRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	m := domain.TeamMember{Name: req.Name, Email: req.Email, Role: req.Role, Status: req.Status}
	if err := h.Team.Create(r.Context(), &m, req.Password); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, teamResponse(m))
}
func (h Handlers) GetTeamMember(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	m, err := h.Team.GetByID(r.Context(), rid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, teamResponse(*m))
}
func (h Handlers) ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	pf := page(r)
	res, err := h.Team.List(r.Context(), repositories.TeamFilter{PageFilter: pf, Role: r.URL.Query().Get("role"), Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	items := make([]TeamResponse, 0, len(res.Items))
	for _, m := range res.Items {
		items = append(items, teamResponse(m))
	}
	writeList(w, repositories.ListResult[TeamResponse]{Items: items, Total: res.Total}, pf.Page, pf.PageSize)
}
func (h Handlers) UpdateTeamMember(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var req TeamRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, domain.ErrValidation)
		return
	}
	m := domain.TeamMember{ID: rid, Name: req.Name, Email: req.Email, Role: req.Role, Status: req.Status}
	if err := h.Team.Update(r.Context(), &m); err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, teamResponse(m))
}
func (h Handlers) DeleteTeamMember(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Team.Delete(r.Context(), rid); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handlers) DeactivateTeamMember(w http.ResponseWriter, r *http.Request) {
	rid, err := id(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Team.Deactivate(r.Context(), rid, actorID(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
