package controllers

import "github.com/steevepypo/todoback/src/services/middlewares"

func (s *Server) Routes() {
	// s.Router.LoadHTMLGlob("templates/*")
	// https://community.fly.io/t/how-to-deploy-static-templates-with-go-golang/1339
	// s.Router.LoadHTMLFiles(filepath.Join("src", "public", "templates", "**/*"))
	// s.Router.LoadHTMLGlob(filepath.Join("src", "public", "templates", "**", "*"))
	// html := template.Must(template.ParseFiles(filepath.Join("src", "public", "templates", "*.html"), filepath.Join("src", "public", "templates", "blocos", "*.html")))
	// s.Router.SetHTMLTemplate(html)

	appUrl := s.Router.Group("/api")
	appUrl.POST("/login", s.Login)
	appUrl.POST("/register", s.CreateUser)
	appUrl.GET("/usuarios", middlewares.CookiesSessionMiddleware(s.DB), s.ListarUsuario)
	appUrl.DELETE("/remover", s.RemoverUsuarios)

	// s.Router.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(
	// 		// Set the HTTP status to 200 (OK)
	// 		http.StatusOK,
	// 		// Use the index.html template
	// 		"bloco.html",
	// 		// Pass the data that the page uses (in this case, 'title')
	// 		gin.H{
	// 			"title": "Steeve home title",
	// 		},
	// 	)
	// })

}
