package v1

type ClientRole string

const (
	adminRole     ClientRole = "admin"
	moderatorRole ClientRole = "moderator"
	supportRole   ClientRole = "support"
)
