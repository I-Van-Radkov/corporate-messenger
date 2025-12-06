package v1

import "github.com/gin-gonic/gin"

type DirectoryUsecase interface {
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

}

func (h *DirectoryHandlers) GetDepartments(c *gin.Context) {

}

func (h *DirectoryHandlers) RemoveDepartment(c *gin.Context) {

}

func (h *DirectoryHandlers) GetDepartmentMembers(c *gin.Context) {

}

func (h *DirectoryHandlers) CreateUser(c *gin.Context) {

}

func (h *DirectoryHandlers) GetUsers(c *gin.Context) {

}

func (h *DirectoryHandlers) GetUser(c *gin.Context) {

}

func (h *DirectoryHandlers) RemoveUser(c *gin.Context) {

}
