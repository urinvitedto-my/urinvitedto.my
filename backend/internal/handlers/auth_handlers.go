package handlers

import "net/http"

// GetMe returns the authenticated user's email and admin status.
func (h *Handlers) GetMe(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || email == "" {
		h.writeError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	isAdmin, err := h.IsAdmin(r.Context(), email)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "server_error", "Failed to check admin status")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"email":   email,
		"isAdmin": isAdmin,
	})
}
