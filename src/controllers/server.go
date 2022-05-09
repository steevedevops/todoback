package controllers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/steevepypo/todoback/src/models"
	"github.com/steevepypo/todoback/src/services/security"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) InitDB() {
	var err error
	dbdriver := os.Getenv("DB_DRIVER")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, db_port, username, database, password)
	server.DB, err = gorm.Open(dbdriver, dsn)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", dbdriver)
		log.Fatal("This is the error connecting to postgres:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", dbdriver)
	}
	AutoMigrate(server)
}

func AutoMigrate(server *Server) {
	server.DB.Debug().AutoMigrate(
		&models.User{},
		&models.Session{},
	)
	server.DB.Model(&models.Session{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
}

func (server *Server) InitRouter() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	server.Router = gin.Default()
	server.Routes()
	log.Fatal(server.Router.Run("0.0.0.0:" + os.Getenv("SV_PORT")))
}

func (server *Server) Authenticate(ctx *gin.Context, username, password string) (*models.User, string, error) {
	user := models.User{}
	sessions := models.Session{}

	var err error
	err = server.DB.Debug().Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &user, "", fmt.Errorf("%s", "Usuário ou Senha invalida!")
	}

	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &user, "", fmt.Errorf("%s", "Usuário ou Senha invalida!")
	}

	/* Continuar com o processo de criar a sessao do usuario */
	sessionid, err := ctx.Cookie("sessionid")

	sessionExp := time.Now().Add(120 * time.Second)
	sessions.UserID = user.UserID
	// security.TokenHashmd5(uuid.NewV4().String())
	user_session_key := fmt.Sprintf("%d:%s:%s:%s:%s", time.Now().Nanosecond(), strconv.Itoa(user.UserID), username, user.FirstName, user.LastName)
	sessionKey, err_encr := security.Encrypt(user_session_key)
	if err_encr != nil {
		fmt.Println("Error encrypting your classified text: ", err)
	}
	sessions.SessionKey = sessionKey
	sessions.UserAgent = ctx.Request.UserAgent()
	sessions.ClientIp = ctx.ClientIP()
	sessions.IsBlocked = false
	sessions.ExpiresAt = sessionExp

	if err != nil {
		session_db, err_db := sessions.Save(server.DB)
		if err_db != nil {
			return &user, "", fmt.Errorf("%s", "Não foi possivel criar sua sessão!")
		}
		ctx.SetCookie("sessionid", session_db.SessionKey, session_db.ExpiresAt.Second(), "/", ctx.Request.Host, false, true)
	} else {
		session_db, err := sessions.FindSessionkey(server.DB, sessionid, user.UserID)
		if err != nil {
			session_db, err_db := sessions.Save(server.DB)
			if err_db != nil {
				return &user, "", fmt.Errorf("%s", "Não foi possivel criar sua sessão!")
			}
			ctx.SetCookie("sessionid", session_db.SessionKey, session_db.ExpiresAt.Second(), "/", ctx.Request.Host, false, true)
		} else {
			ctx.SetCookie("sessionid", session_db.SessionKey, session_db.ExpiresAt.Second(), "/", ctx.Request.Host, false, true)
		}
	}
	return &user, sessions.SessionKey, nil
}
