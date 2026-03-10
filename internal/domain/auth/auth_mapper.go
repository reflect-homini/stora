package auth

import (
	"github.com/reflect-homini/stora/internal/domain/appconstant"
)

func sessionToAuthData(session Session) map[string]any {
	return map[string]any{
		string(appconstant.ContextUserID):    session.UserID,
		string(appconstant.ContextSessionID): session.ID,
	}
}
