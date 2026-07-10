package server

import (
	"github.com/gin-gonic/gin"
	"github.com/shapnesht/ecommerce-platform-go/internal/dto"
	"github.com/shapnesht/ecommerce-platform-go/internal/utils"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid Request data", err)
		return
	}

	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid Request data", err)
		return
	}

	response, err := s.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login failed")
		return
	}

	utils.SuccessResponse(c, "Logged in successfully", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid format", err)
	}

	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "refresh token failed")
		return
	}

	utils.SuccessResponse(c, "token refresh successful", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid format", err)
	}

	err := s.authService.Logout(req.RefreshToken)
	if err != nil {
		utils.InternalServerErrorResponse(c, "logout failed", err)
		return
	}

	utils.SuccessResponse(c, "logout successful", nil)
}

func (s *Server) getProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userId)

	if err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	utils.SuccessResponse(c, "user found successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	var req *dto.ProfileUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid update request data", err)
	}

	profile, err := s.userService.UpdateProfile(userId, req)

	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "user updated successfully", profile)
}
