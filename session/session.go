package session

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

type SessionHandler struct {
	sc *securecookie.SecureCookie
}

func NewSessionHandler(hashKeyLen, blockKeyLen int) *SessionHandler {
	return &SessionHandler{
		sc: securecookie.New(
			securecookie.GenerateRandomKey(hashKeyLen),
			securecookie.GenerateRandomKey(blockKeyLen),
		),
	}
}

func (s *SessionHandler) GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = s.sc.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}

func (s *SessionHandler) SetSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}
	if encoded, err := s.sc.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func (s *SessionHandler) ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
