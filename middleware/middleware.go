package middleware

import (
    "net/http"
    "net/http/httputil"
    "fmt"
    "context"
    "encoding/json"
)

type AuthResponse struct {
    UserId string `json:"user_id"`
}

func AuthMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
        client := &http.Client{}
        req, err := http.NewRequest("POST", "http://localhost/api/v1/users/auth", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        req.Header.Set("Authorization", r.Header.Get("Authorization"))
        resp, err := client.Do(req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        defer resp.Body.Close()

        var a AuthResponse
        err = json.NewDecoder(resp.Body).Decode(&a)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        newCtx := context.WithValue(r.Context(), "uid", a.UserId)
        next.ServeHTTP(w, r.WithContext(newCtx))
    })
}

func RequestDump(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(dump))
		next.ServeHTTP(w, r)
    })
}
