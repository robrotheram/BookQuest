package auth

import (
	"BookQuest/internal/models"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
)

type UserReqContext string

const userCtx = UserReqContext("user")

func GetUser(r *http.Request) (models.User, error) {
	user, ok := r.Context().Value(userCtx).(models.User)
	if !ok {
		return user, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func UpdateUser(db *bun.DB, user models.User) {

}

func AuthMiddleware(db *bun.DB, authN *authentication.Authenticator[*openid.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]) func(next http.Handler) http.Handler {

	userCache := initCache()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := authN.IsAuthenticated(r)
			if err != nil {
				authN.Authenticate(w, r, r.URL.Path)
				return
			}
			user, err := userCache.GetUser(session.UserInfo.Email)
			if err != nil {
				user, err = models.GetUserByEmail(db, session.UserInfo.Email)
				if err != nil {
					user = models.User{
						Id:       uuid.New(),
						Email:    session.UserInfo.Email,
						Username: session.UserInfo.Name,
						Picture:  session.UserInfo.Picture,
					}
					models.CreateUser(db, user)
				}
				user.Picture = session.UserInfo.Picture
				if len(user.Picture) == 0 {
					user.Picture = "/icon/" + session.UserInfo.Email
				}
				models.UpdateUser(db, user)
				userCache.SetUser(user)
			}
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}

			//Update User
			go models.UpdateUser(db, user)
			go userCache.SetUser(user)
			ctx := context.WithValue(r.Context(), userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
