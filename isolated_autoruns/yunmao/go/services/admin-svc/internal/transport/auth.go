package transport

import (
	"context"
	"net/http"
	"strings"

	"yunmao.live/pkg/yunmao/authjwt"
	yerr "yunmao.live/pkg/yunmao/errors"
	"yunmao.live/pkg/yunmao/httpx"
)

type claimsKey struct{}

func RequireAdmin(verifier *authjwt.Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if verifier == nil {
				next.ServeHTTP(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				httpx.WriteError(w, r, yerr.New(yerr.AuthLoginRequired, "missing bearer token"))
				return
			}
			claims, err := verifier.Parse(strings.TrimPrefix(auth, "Bearer "))
			if err != nil {
				httpx.WriteError(w, r, yerr.New(yerr.AuthTokenExpired, err.Error()))
				return
			}
			if claims.Scope != authjwt.ScopeAdmin {
				httpx.WriteError(w, r, yerr.New(yerr.AuthForbidden, "admin scope required"))
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminClaims(ctx context.Context) *authjwt.Claims {
	v, _ := ctx.Value(claimsKey{}).(*authjwt.Claims)
	return v
}

func AdminClaimsFromRouter(r *http.Request) *authjwt.Claims {
	return AdminClaims(r.Context())
}
