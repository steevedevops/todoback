package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/steevepypo/todoback/src/models"
	"github.com/steevepypo/todoback/src/services/security"
)

func CookiesSessionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessions := models.Session{}
		user := models.User{}
		ctx.Set("user_authenticated", false)
		ctx.Set("user", models.User{})
		msg := "As credenciais de autenticação não foram fornecidas."
		sessionid_cookie, err := ctx.Cookie("sessionid")
		if err != nil {
			ctx.Set("user", models.User{})
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": msg,
			})
			ctx.Abort()
			return
		}
		sessionId, err := security.Decrypt(sessionid_cookie)
		if err != nil {
			ctx.Set("user", models.User{})
			fmt.Println("Error encrypting your classified text: ", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user_session := strings.Split(sessionId, ":")
		session_db := sessions.FindSessionUser(db, user_session[1], sessionid_cookie)
		if !session_db && sessions.IsExpired() {
			sessions.Delete(db, sessionid_cookie)
			ctx.Set("user", models.User{})
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": msg,
			})
			ctx.Abort()
			return
		}
		if session_db {
			ctx.Set("user", models.User{})
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": msg,
			})
			ctx.Abort()
			return
		}
		user_db, err := user.FindById(db, user_session[1])
		if err != nil {
			ctx.Set("user_authenticated", false)
			ctx.Set("user", models.User{})
		} else {
			ctx.Set("user_authenticated", true)
			ctx.Set("user", user_db)
		}
		ctx.Next()

	}
}
