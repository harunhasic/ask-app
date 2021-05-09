package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mop/config"
	"github.com/mop/controller"
	"github.com/mop/middleware"
	"github.com/mop/repository"
	"github.com/mop/service"
	"gorm.io/gorm"
)

var (
	db                 *gorm.DB                      = config.SetupDatabaseConnection()
	userRepository     repository.UserRepository     = repository.NewUserRepository(db)
	questionRepository repository.QuestionRepository = repository.NewQuestionRepository(db)
	answerRepository   repository.AnswerRepository   = repository.NewAnswerRepository(db)
	userService        service.UserService           = service.NewUserService(userRepository)
	questionService    service.QuestionService       = service.NewQuestionService(questionRepository)
	answerService      service.AnswerService         = service.NewAnswerService(answerRepository)
	jwtService         service.JWTService            = service.NewJWTService()
	authService        service.AuthService           = service.NewAuthService(userRepository)
	authController     controller.AuthController     = controller.NewAuthController(authService, jwtService)
	userController     controller.UserController     = controller.NewUserController(userService, jwtService)
	answerController   controller.AnswerController   = controller.NewAnswerController(answerService, jwtService)
	questionController controller.QuestionController = controller.NewQuestionController(questionService, jwtService)
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	r := gin.New()
	r.Use(CORSMiddleware())
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile/:id", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	mainPageRoutes := r.Group("api/main")
	{
		mainPageRoutes.GET("/", questionController.All)
		mainPageRoutes.GET("/likes", questionController.GetNumberOfLikesForQuestion)
		mainPageRoutes.GET("/answers", userController.GetAnswers)
		mainPageRoutes.GET("/question:id", questionController.QuestionPage)
	}

	questionRoutes := r.Group("api/questions", middleware.AuthorizeJWT(jwtService))
	{
		questionRoutes.POST("/like/:id", questionController.Like)
		questionRoutes.DELETE("/like/:id", questionController.DeleteLike)
		questionRoutes.GET("/likes", questionController.GetNumberOfLikesForQuestion)
		questionRoutes.POST("/", questionController.Insert)
		questionRoutes.GET("/:id", questionController.FindByID)
		questionRoutes.GET("/page/:id", questionController.QuestionPage)
		questionRoutes.PUT("/", questionController.Update)
		questionRoutes.DELETE("/:id", questionController.Delete)
		questionRoutes.GET("/profile/:id", questionController.GetUserQuestions)
	}

	answerRoutes := r.Group("api/answers", middleware.AuthorizeJWT(jwtService))
	{
		answerRoutes.GET("/", answerController.All)
		answerRoutes.POST("/", answerController.Insert)
		answerRoutes.GET("/:id", answerController.FindByID)
		answerRoutes.PUT("/", answerController.Update)
		answerRoutes.DELETE("/:id", answerController.Delete)
	}

	r.Run()
}
