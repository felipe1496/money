package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"rango-backend/resources/users"
	"rango-backend/services"
	"rango-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type API struct {
	googleService services.GoogleService
	usersUseCase  users.UsersUseCase
	JWTService    services.JWTService
}

func NewHandler(db *sql.DB, googleService services.GoogleService, jwtService services.JWTService) *API {
	return &API{
		googleService: googleService,
		usersUseCase:  users.NewUsersUseCase(users.NewUsersRepo(db)),
		JWTService:    jwtService,
	}
}

func (h *API) LoginGoogle(ctx *gin.Context) {
	var body LoginGoogleRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "It was not possible to process the request body"))
		return
	}

	userAccessToken, err := h.googleService.GetUserAccessToken(body.Code)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.ErrorResponse(http.StatusInternalServerError, "Failed to get user access token with Google"))
		return
	}

	userInfo, err := h.googleService.GetUserInfo(*userAccessToken)

	if err != nil {
		if errors.Is(err, GoogleAuthFailedErr) {
			ctx.JSON(http.StatusUnauthorized,
				utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError,
				utils.ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	userExists, err := h.usersUseCase.List(users.UserFilter{Email: *userInfo.Email})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch user information"))
		return
	}

	var userRes users.User
	var status int

	if len(userExists) == 0 {
		createUserInput := users.CreateUserInput{
			Name: userInfo.Name,
		}

		if userInfo.Email == nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Email not provided by Google"))
			return
		}

		createUserInput.Email = *userInfo.Email

		createUserInput.AvatarURL = userInfo.Picture

		createUserInput.Username = fmt.Sprintf("%s_%s", strings.ToLower(strings.ReplaceAll(userInfo.Name, " ", "_")), ulid.Make().String())

		createdUser, err := h.usersUseCase.Create(createUserInput)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create user"))
			return
		}

		userRes = createdUser

		status = http.StatusCreated
	} else {
		userRes = userExists[0]

		status = http.StatusOK
	}

	access_token, err := h.JWTService.GenerateToken(userRes.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to generate access token"))
		return
	}

	ctx.JSON(status, gin.H{
		"user":         userRes,
		"access_token": access_token,
	})
}
