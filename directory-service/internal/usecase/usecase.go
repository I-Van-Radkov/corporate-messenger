package usecase

type DirectoryRepo interface {
}

type DirectoryUsecase struct {
	Repository DirectoryRepo
}

func NewDirectoryService(repository DirectoryRepo) *DirectoryUsecase {
	return &DirectoryUsecase{
		Repository: repository,
	}
}
