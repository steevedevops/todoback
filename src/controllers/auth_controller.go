package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steevepypo/todoback/src/models"
)

func (server *Server) Login(ctx *gin.Context) {
	var urs models.User

	// Verifica se o payload veio certo
	if err := ctx.ShouldBindJSON(&urs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, sessionId, err := server.Authenticate(ctx, urs.Username, urs.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message":    "Sucesso ao fazer login",
		"user":       user,
		"session_id": sessionId,
	})
}

func (db *Server) CreateUser(ctx *gin.Context) {
	var u models.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userCreated, err := u.Save(db.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "usuario creado com sucesso!",
		"user":    userCreated,
	})
}

func (db *Server) ListarUsuario(ctx *gin.Context) {
	user := models.User{}
	user_logged := ctx.MustGet("user").(*models.User)
	users, err := user.List(db.DB, user_logged.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (db *Server) RemoverUsuarios(ctx *gin.Context) {
	user := models.User{}

	err := user.Delete(db.DB)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"users": "Usuarios deletado com sucesso",
	})
}
