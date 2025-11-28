package v1

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
