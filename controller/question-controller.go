package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mop/dto"
	"github.com/mop/entity"
	"github.com/mop/helper"
	"github.com/mop/service"
)

//QuestionController is the interface that shows what methods a question
//controller needs to have/implement
type QuestionController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	GetNumberOfLikesForQuestion(context *gin.Context)
	Like(context *gin.Context)
	DeleteLike(context *gin.Context)
	QuestionPage(context *gin.Context)
}

type questionController struct {
	questionService service.QuestionService
	jwtService      service.JWTService
}

//NewQuestionController creates a new instance of the QuestionController
func NewQuestionController(questionServ service.QuestionService, jwtServ service.JWTService) QuestionController {
	return &questionController{
		questionService: questionServ,
		jwtService:      jwtServ,
	}
}

//All returns all the question entries in the db
func (c *questionController) All(context *gin.Context) {
	var questions []entity.Question = c.questionService.All()
	res := helper.BuildResponse(true, "OK!", questions)
	context.JSON(http.StatusOK, res)
}

func (c *questionController) GetNumberOfLikesForQuestion(context *gin.Context) {
	var questions []entity.Question = c.questionService.GetNumberOfLikesForQuestion()
	res := helper.BuildResponse(true, "OK!", questions)
	context.JSON(http.StatusOK, res)
}

//FindByID method returns a question with the given ID
//Returns the question if there is a match
//If a question with no such ID exists, returns Data not found response
func (c *questionController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var question entity.Question = c.questionService.FindById(id)
	if (question == entity.Question{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK!", question)
		context.JSON(http.StatusOK, res)
	}
}

func (c *questionController) QuestionPage(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	theId, errID := strconv.ParseUint(userID, 10, 64)
	if errID == nil {
		var questionResult entity.QuestionResult = c.questionService.QuestionPage(id, theId)
		response := helper.BuildResponse(true, "OK", questionResult)
		context.JSON(http.StatusOK, response)
	} else {
		res := helper.BuildErrorResponse("No id found", "No data with that id", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

}

func (c *questionController) Insert(context *gin.Context) {
	var questionCreateDTO dto.QuestionCreateDTO
	errDTO := context.ShouldBind(&questionCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process the request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			questionCreateDTO.UserID = convertedUserID
		}
		result := c.questionService.Insert(questionCreateDTO)
		response := helper.BuildResponse(true, "OK!", result)
		context.JSON(http.StatusCreated, response)
	}
}

//The update method is used mainly to update/edit the question that a user posts.
//The method incorporates the IsAllowedToEdit method which checks if the correct user is trying to edit the question
func (c *questionController) Update(context *gin.Context) {
	var questionUpdateDTO dto.QuestionUpdateDTO
	errDTO := context.ShouldBind(&questionUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.questionService.IsAllowedToEdit(userID, questionUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			questionUpdateDTO.UserID = id
		}
		result := c.questionService.Update(questionUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have the permission", "You are not the author of the question", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

//Delete method is used for the deletion of questions
func (c *questionController) Delete(context *gin.Context) {
	var question entity.Question
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get ID", "No param id was found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	question.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.questionService.IsAllowedToEdit(userID, question.ID) {
		c.questionService.Delete(question)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have the permission", "You are not the author of the question", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *questionController) Like(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get question id", "No param question id was found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	theId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		panic(err)
	}
	c.questionService.Like(id, uint64(theId))
	response := helper.BuildResponse(true, "Liked", helper.EmptyObj{})
	context.JSON(http.StatusOK, response)
}

func (c *questionController) DeleteLike(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get question id", "No param question id was found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	theId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		panic(err)
	}
	c.questionService.DeleteLike(id, uint64(theId))
	response := helper.BuildResponse(true, "Unliked", helper.EmptyObj{})
	context.JSON(http.StatusOK, response)
}

//getUserIDByToken retrieves the user id from the given token
func (c *questionController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
