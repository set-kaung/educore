package auth

type SetupContext struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SetupRequest struct {
	Name         string `json:"name"`
	DepartmentID uint   `json:"department_id"`
	StudentID    string `json:"student_id"`
}

type SetupResult struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
