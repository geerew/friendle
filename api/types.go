package api

// versionResponse is the application version payload
type versionResponse struct {
	Version string `json:"version"`
}

// signupStatusResponse reports whether self-service registration is enabled
type signupStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// adminRecoveryRequest is the body for POST /api/admin/recovery
type adminRecoveryRequest struct {
	Token string `json:"token"`
}
