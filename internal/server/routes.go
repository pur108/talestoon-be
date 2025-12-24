package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/pur108/webteen-be/internal/delivery/http"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/middleware"
	"github.com/pur108/webteen-be/internal/repository"
	"github.com/pur108/webteen-be/internal/usecase"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(logger.New())
	s.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Get("/", s.HelloWorldHandler)
	s.Get("/health", s.healthHandler)

	db := s.db.GetDB()

	// user routes
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(s.App, userUsecase)

	// auth routes
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := http.NewAuthHandler(s.App, authUsecase)

	// comic routes
	comicRepo := repository.NewComicRepository(db)
	storageRepo := repository.NewSupabaseRepository()
	comicUsecase := usecase.NewComicUsecase(comicRepo, userRepo, storageRepo)
	comicHandler := http.NewComicHandler(s.App, comicUsecase)

	uploadUsecase := usecase.NewUploadUsecase(storageRepo)
	uploadHandler := http.NewUploadHandler(s.App, uploadUsecase)

	// Admin Routes
	adminHandler := http.NewAdminHandler(s.App, comicUsecase, storageRepo)

	group := s.Group("/api")

	// Admin Group
	adminGroup := group.Group("/admin", middleware.Protected(), middleware.RoleRequired(entity.RoleAdmin))
	adminGroup.Get("/comics", adminHandler.ListPendingComics)
	adminGroup.Post("/comics/:id/approve", adminHandler.ApproveComic)
	adminGroup.Post("/comics/:id/reject", adminHandler.RejectComic)

	// Auth Routes
	group.Post("/auth/signup", authHandler.SignUp)
	group.Post("/auth/login", authHandler.Login)

	// User Routes
	userGroup := group.Group("/users", middleware.Protected())
	userGroup.Get("/me", userHandler.GetProfile)
	userGroup.Post("/become-creator", userHandler.BecomeCreator)

	// Comic Public Routes
	group.Get("/comics", comicHandler.ListComics)
	group.Get("/comics/:id", comicHandler.GetComic)
	group.Get("/chapters/:id", comicHandler.GetChapter)
	group.Get("/tags", comicHandler.ListTags)

	// Comic Creator Routes
	creatorGroup := group.Group("/creator/comics", middleware.Protected(), middleware.RoleRequired(entity.RoleCreator, entity.RoleAdmin, entity.RoleUser))
	creatorGroup.Post("", comicHandler.CreateComic)
	creatorGroup.Get("", comicHandler.ListMyComics)
	creatorGroup.Put("/:id", comicHandler.UpdateComic)
	creatorGroup.Delete("/:id", comicHandler.DeleteComic)
	creatorGroup.Post("/:id/chapters", comicHandler.CreateChapter)
	creatorGroup.Post("/:id/publish-request", comicHandler.RequestPublish)

	// Library Routes
	// libGroup := group.Group("/library")

	// // Public Library Folders
	// libGroup.Get("/public/folders/:slug", libraryHandler.GetPublicFolder)

	// Protected Library Routes
	// libProtected := libGroup.Group("/", middleware.Protected())
	// libProtected.Post("/entries", libraryHandler.AddToLibrary)
	// libProtected.Delete("/entries/:comic_id", libraryHandler.RemoveFromLibrary)
	// libProtected.Get("/entries", libraryHandler.GetUserLibrary)

	// Protected Folder Management
	// libProtected.Post("/folders", libraryHandler.CreateFolder)
	// libProtected.Delete("/folders/:id", libraryHandler.DeleteFolder)
	// libProtected.Get("/folders", libraryHandler.GetUserFolders)
	// libProtected.Get("/folders/:id", libraryHandler.GetFolder)

	// Protected Folder Items
	// libProtected.Post("/folders/:id/items", libraryHandler.AddToFolder)
	// libProtected.Delete("/folders/:id/items/:comic_id", libraryHandler.RemoveFromFolder)

	// Upload Routes
	group.Post("/upload", middleware.Protected(), uploadHandler.UploadFile)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
