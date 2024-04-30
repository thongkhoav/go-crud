package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thongkhoav/go-crud/controllers"
	"github.com/thongkhoav/go-crud/initializers"
	"github.com/thongkhoav/go-crud/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {

	routes := gin.Default()

	routes.POST("/signup", controllers.Signup)
	routes.POST("/login", controllers.Login)
	routes.GET("/protect", middleware.RequireAuth, controllers.TestAuthenication)

	routes.POST("/posts", middleware.RequireAuth, controllers.PostCreate)
	routes.PUT("/posts/:id", middleware.RequireAuth, controllers.PostUpdate)
	routes.DELETE("/posts/:id", middleware.RequireAuth, controllers.PostDelete)
	routes.GET("/posts", controllers.PostIndex)
	routes.GET("/posts/:id", controllers.PostShow)

	routes.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
