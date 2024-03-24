package middlewares

import "net/http"

const (
	responseHeaderSetCookie = "Set-Cookie"
)

type hardenAuthCookieResponseWriter struct {
	authCookieName     string
	baseResponseWriter http.ResponseWriter
}

func NewHardenAuthCookieResponseWriter(
	authCookieName string,
	baseResponseWriter http.ResponseWriter,
) http.ResponseWriter {
	return &hardenAuthCookieResponseWriter{
		authCookieName:     authCookieName,
		baseResponseWriter: baseResponseWriter,
	}
}

func (h hardenAuthCookieResponseWriter) Header() http.Header {
	return h.baseResponseWriter.Header()
}

func (h hardenAuthCookieResponseWriter) getCookies(cookieHeaderValues []string) []*http.Cookie {
	return http.Request{
		Header: http.Header{
			responseHeaderSetCookie: cookieHeaderValues,
		},
	}.Response.Cookies()
}

func (h hardenAuthCookieResponseWriter) WriteHeader(statusCode int) {
	setCookieHeaderValues := h.baseResponseWriter.Header().Values(responseHeaderSetCookie)
	if len(setCookieHeaderValues) == 0 {
		h.baseResponseWriter.WriteHeader(statusCode)
		return
	}

	h.baseResponseWriter.Header().Del(responseHeaderSetCookie)

	cookies := h.getCookies(setCookieHeaderValues)
	for _, cookie := range cookies {
		if cookie.Name != h.authCookieName {
			http.SetCookie(h.baseResponseWriter, cookie)
			continue
		}

		cookie.SameSite = http.SameSiteStrictMode
		cookie.HttpOnly = true
		http.SetCookie(h.baseResponseWriter, cookie)
	}

	h.baseResponseWriter.WriteHeader(statusCode)
}

func (h hardenAuthCookieResponseWriter) Write(data []byte) (int, error) {
	return h.baseResponseWriter.Write(data)
}

type HardenAuthCookie func(http.Handler) http.Handler

func NewHardenAuthCookie(authCookieName string) HardenAuthCookie {
	return func(baseHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			baseHandler.ServeHTTP(NewHardenAuthCookieResponseWriter(authCookieName, w), r)
		})
	}
}
