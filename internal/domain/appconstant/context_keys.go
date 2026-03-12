package appconstant

type ctxKey string

const (
	ContextUserID    ctxKey = "userID"
	ContextSessionID ctxKey = "sessionID"

	ContextProvider ctxKey = "provider"

	ContextProjectID ctxKey = "projectID"
)
