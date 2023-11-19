package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/app/dto"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/app/service"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/config"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/helpers"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/models"
	"github.com/hasrulrhul/service-repository-pattern-gin-golang/response"
)

// UserController is a contract what this controller can do
type UserController interface {
	Index(ctx *gin.Context)
	Create(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

// NewUserController create a new instances of UserController
func NewUserController(userServ service.UserService, jwtServ service.JWTService) UserController {
	return &userController{
		userService: userServ,
		jwtService:  jwtServ,
	}
}

func (s *userController) Index(ctx *gin.Context) {
	pagination := helpers.GeneratePaginationRequest(ctx)
	res := s.userService.Pagination(ctx, pagination)
	if !res.Status {
		response := response.ResponseError("failed to get data user", res.Message)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := response.ResponseSuccess("list of user", res.Data)
	ctx.JSON(http.StatusOK, response)
}

func (s *userController) Create(ctx *gin.Context) {
	var req dto.UserCreateValidation
	err := ctx.ShouldBind(&req)
	if err != nil {
		response := response.ResponseError(config.MessageErr.FailedProcess, err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// if !s.userService.FindByEmail(req.Email) {
	// 	response := response.ResponseError(config.MessageErr.FailedProcess, "duplicate email")
	// 	ctx.JSON(http.StatusConflict, response)
	// } else {
	err = s.userService.Create(req)
	if err != nil {
		go helpers.CreateLogError(uuid.NewString(), helpers.GetIP(ctx), "users", "created users", err.Error())
		response := response.ResponseError("failed to process created", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	go helpers.CreateLogInfo(uuid.NewString(), helpers.GetIP(ctx), "users", "created users", "created success")
	response := response.ResultSuccess("created success")
	ctx.JSON(http.StatusCreated, response)
	// }
}

func (s *userController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User = s.userService.Show(id)
	if (user == models.User{}) {
		res := response.ResponseError("Data not found", "No data with given id")
		ctx.JSON(http.StatusNotFound, res)
	} else {
		response := response.ResponseSuccess("detail user", user)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *userController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User = c.userService.Show(id)
	if user.ID == "" {
		res := response.ResponseError("data not found", "no data with given id")
		ctx.JSON(http.StatusNotFound, res)
	} else {
		var userValidation dto.UserUpdateValidation
		userValidation.ID = id
		err := ctx.ShouldBind(&userValidation)
		if err != nil {
			response := response.ResponseError(config.MessageErr.FailedProcess, err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		err = c.userService.Update(userValidation)
		if err != nil {
			go helpers.CreateLogError(uuid.NewString(), helpers.GetIP(ctx), "users", "updated users", err.Error())
			response := response.ResponseError("update failed", err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		go helpers.CreateLogInfo(uuid.NewString(), helpers.GetIP(ctx), "users", "updated users", "updated success")
		response := response.ResponseSuccess("update success", nil)
		ctx.JSON(http.StatusCreated, response)
	}
}

func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User = c.userService.Show(id)
	if user.ID == "" {
		response := response.ResponseError("data not found", "no data with given id")
		ctx.JSON(http.StatusNotFound, response)
	} else {
		err := c.userService.Delete(user)
		if err != nil {
			go helpers.CreateLogError(uuid.NewString(), helpers.GetIP(ctx), "users", "deleted users", err.Error())
			response := response.ResponseError("failed to process deleted", err.Error())
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		go helpers.CreateLogInfo(uuid.NewString(), helpers.GetIP(ctx), "users", "deleted users", "deleted success")
		response := response.ResultSuccess("deleted success")
		ctx.JSON(http.StatusOK, response)
	}
}
