package schemas

type PreRegistrationInput struct {
	Name        string `json:"name" validate:"required,name" example:"John Doe"`
	Email       string `json:"email" validate:"required,email" example:"example@mail.com"`
	DateOfBirth string `json:"date_of_birth" example:"1990-01-01"`
	PhoneNumber string `json:"phone_number" example:"08123456789"`
	Password    string `json:"password" validate:"required,password" example:"password123"`
}

type VerifyOTPInput struct {
	Email string `json:"email" validate:"required,email" example:"example@mail.com"`
	OTP   string `json:"otp" validate:"required" example:"123456"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email" example:"example@mail.com"`
	Password string `json:"password" validate:"required,password" example:"password123"`
}
