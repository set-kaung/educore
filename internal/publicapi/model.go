package publicapi

type VerifyDepartmentResponse struct {
	StudentID    string `json:"student_id"`
	Department   string `json:"department"`
	Enrolled     bool   `json:"enrolled"`
}

type GrantKeyRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GrantKeyResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Key   string `json:"key"`
}
