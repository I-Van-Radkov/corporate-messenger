package v1

import (
	"context"
	"net/http"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DirectoryUsecase interface {
	CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.CreateDepartmentResponse, error)
	GetDepartments(ctx context.Context, req *dto.GetDepartmentsRequest) (*dto.GetDepartmentsResponse, error)
	RemoveDepartment(ctx context.Context, id uuid.UUID) error
	GetDepartmentMembers(ctx context.Context, req *dto.GetDepartmentMembersRequest) (*dto.GetDepartmentMembersResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetUsers(ctx context.Context, req *dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	RemoveUser(ctx context.Context, id uuid.UUID) error
}

type DirectoryHandlers struct {
	usecase DirectoryUsecase
}

func NewDirectoryHandlers(usecase DirectoryUsecase) *DirectoryHandlers {
	return &DirectoryHandlers{
		usecase: usecase,
	}
}

func (h *DirectoryHandlers) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.CreateDepartment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DirectoryHandlers) GetDepartments(c *gin.Context) {
	var req *dto.GetDepartmentsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.GetDepartments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DirectoryHandlers) RemoveDepartment(c *gin.Context) {
	depIdStr := c.Param("department_id")
	depId, err := uuid.Parse(depIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department_id"})
		return
	}

	err = h.usecase.RemoveDepartment(c.Request.Context(), depId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DirectoryHandlers) GetDepartmentMembers(c *gin.Context) {
	depIdStr := c.Param("department_id")
	depId, err := uuid.Parse(depIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department_id"})
		return
	}

	var req dto.GetDepartmentMembersRequest
	req.DepartmentID = depId
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.GetDepartmentMembers(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DirectoryHandlers) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DirectoryHandlers) GetUsers(c *gin.Context) {
	var req dto.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.GetUsers(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DirectoryHandlers) GetUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	user, err := h.usecase.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *DirectoryHandlers) RemoveUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	err = h.usecase.RemoveUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
