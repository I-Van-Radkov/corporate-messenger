package identity

type IntrospectRequest struct {
	Token string `json:"token"`
}

type IntrospectResponse struct {
	Active bool `json:"active"`
}
