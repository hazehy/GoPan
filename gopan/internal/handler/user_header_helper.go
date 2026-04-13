package handler

import "net/http"

func requestUserIdentity(r *http.Request) string {
	if userIdentity := r.Header.Get("UserIdentity"); userIdentity != "" {
		return userIdentity
	}
	return r.Header.Get("userIdentity")
}
