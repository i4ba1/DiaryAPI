package user_management

type UserService struct {
	userRepository UserRepository
}

func ProvideUserService(repo UserRepository) UserService {
	return UserService{userRepository: repo}
}

func (a *UserService) signIn(dto LoginDto) (Account, error) {
	return a.userRepository.findUserByUserNameOrEmailAndPass(dto)
}

func (a *UserService) findByUsername(username string) error {
	return a.userRepository.findUserByUserName(username)
}
