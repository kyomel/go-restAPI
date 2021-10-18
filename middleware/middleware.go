package middleware

import (
	"go-rest-api/application/models"
	"io"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func Auth(db *gorm.DB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Access from ", r.URL.Path)

			notAuth := []string{
				"/healthcheck",
				"/api/v1/register",
				"/image",
			}

			withAuth := true
			for _, path := range notAuth {
				if path == r.URL.Path || strings.HasPrefix(r.URL.Path, path) {
					withAuth = false
					break
				}
			}
			if withAuth {
				name, password, ok := r.BasicAuth()
				if !ok {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, `{"Error":"Internal Server Error"}`)
					return
				}

				// pengecekan ke dalam database
				var user models.User
				tx := db.Where("name = ? AND password = ?", name, password).First(&user)
				if tx.Error != nil || tx.RowsAffected == 0 {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					io.WriteString(w, `{"Error":"User not Authorized"}`)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
