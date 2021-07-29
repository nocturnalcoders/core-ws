package user

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageURL   string `json:"image_url"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Role:       user.Role,
		Email:      user.Email,
		Token:      token,
		ImageURL:   user.AvatarFileName,
		// Token:      user.Token,
	}

	return formatter

}
