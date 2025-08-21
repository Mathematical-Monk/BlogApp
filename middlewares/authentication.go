package middlewares

import (
	"blogapi/models"
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var claims models.Claims

		jwtToken, err := r.Cookie("token")
		if err != nil {
			writeResponseForFailedAuthentication(w, err)
			return

		}

		token, err := jwt.ParseWithClaims(jwtToken.Value, &claims, getSecretKey)
		if err != nil {
			writeResponseForFailedAuthentication(w, err)
			return
		}

		ok := token.Valid
		if !ok {
			fmt.Println(ok," token is not valid")
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"authorized":"false"}`))
			return

		}

		_, ok = token.Claims.(*models.Claims)
		if !ok {
			fmt.Println(ok, " claims are not valid")
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"authorized":"false"}`))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)

		r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

func getSecretKey(token *jwt.Token) (any, error) {

	return []byte("my-secret-key"), nil

}

func writeResponseForFailedAuthentication(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"authorized":"false"}`))
}
