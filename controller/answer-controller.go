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

//AnswerController is the interface that shows what methods a answer
//controller needs to have/implement
type AnswerController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	UsersWithMostAnswers(context *gin.Context)
}

type answerController struct {
	answerService service.AnswerService
	jwtService    service.JWTService
}

//NewAnswerController creates a new instance of the AnswerController
func NewAnswerController(answerServ service.AnswerService, jwtServ service.JWTService) AnswerController {
	return &answerController{
		answerService: answerServ,
		jwtService:    jwtServ,
	}
}

//All returns all the answer entries in the db
func (c *answerController) All(context *gin.Context) {
	var answers []entity.Answer = c.answerService.All()
	res := helper.BuildResponse(true, "OK!", answers)
	context.JSON(http.StatusOK, res)
}

func (c *answerController) UsersWithMostAnswers(context *gin.Context) {
	var answers []entity.Answer = c.answerService.UsersWithMostAnswers()
	res := helper.BuildResponse(true, "OK", answers)
	context.JSON(http.StatusOK, res)
}

//FindByID method returns a answer with the given ID
//Returns the answer if there is a match
//If a answer with no such ID exists, returns Data not found response
func (c *answerController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var answer entity.Answer = c.answerService.FindById(id)
	if (answer == entity.Answer{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK!", answer)
		context.JSON(http.StatusOK, res)
	}
}

func (c *answerController) Insert(context *gin.Context) {
	var answerCreateDTO dto.AnswerCreateDTO
	errDTO := context.ShouldBind(&answerCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process the request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			answerCreateDTO.UserID = convertedUserID
		}
		result := c.answerService.Insert(answerCreateDTO)
		response := helper.BuildResponse(true, "OK!", result)
		context.JSON(http.StatusCreated, response)
	}
}

//The update method is used mainly to update/edit the answer that a user posts.
//The method incorporates the IsAllowedToEdit method which makes sure that the user that is trying to edit the answer is the author
func (c *answerController) Update(context *gin.Context) {
	var answerUpdateDTO dto.AnswerUpdateDTO
	errDTO := context.ShouldBind(&answerUpdateDTO)
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
	if c.answerService.IsAllowedToEdit(userID, answerUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			answerUpdateDTO.UserID = id
		}
		result := c.answerService.Update(answerUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have the permission", "You are not the author of the answer", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

//Delete method is used for the deletion of answers
func (c *answerController) Delete(context *gin.Context) {
	var answer entity.Answer
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get ID", "No param id was found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	answer.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.answerService.IsAllowedToEdit(userID, answer.ID) {
		c.answerService.Delete(answer)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have the permission", "You are not the author of the answer", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

//getUserIDByToken retrieves the user id from the given token
func (c *answerController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
