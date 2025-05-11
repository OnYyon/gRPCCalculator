package api

import (
	"net/http"
	"strings"

	"github.com/OnYyon/gRPCCalculator/internal/transport/grpc/auth"
)

type AuthHandler struct {
	authInterceptor *auth.AuthGRPC
	mux             http.Handler
	publicEndpoints []string
}

func NewAuthHandler(authInterceptor *auth.AuthGRPC, mux http.Handler, publicEndpoints []string) *AuthHandler {
	return &AuthHandler{
		authInterceptor: authInterceptor,
		mux:             mux,
		publicEndpoints: publicEndpoints,
	}
}

func (h *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, endpoint := range h.publicEndpoints {
			if r.URL.Path == endpoint {
				h.mux.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Не был префикс Bearer
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		_, err := h.authInterceptor.ValidateTokenAndGetUserID(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		h.mux.ServeHTTP(w, r)
	})
}
