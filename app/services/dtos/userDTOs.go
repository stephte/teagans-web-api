package dtos


type UserDTO struct {
	BaseDTO
	FirstName 	string	`json:"firstName"`
	LastName		string	`json:"lastName"`
	Email				string	`json:"email"`
	Role				int			`json:"role"`
}

// How to get password so we can retrieve it, but not send it??
type CreateUserDTO struct {
	FirstName 	string
	LastName		string
	Email				string
	Role				int
	Password		string
}
