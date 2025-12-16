package model

type Student struct {
	StudentID    string `json:"id"`
	UserID       string `json:"user_id"`
	AcademicYear string `json:"academic_year"`
	ProgramStudy string `json:"program_study"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Fullname     string `json:"full_name"`
	RoleID       string `json:"role_id"` // digunakan agar ketika user menambahkan role mahasiswa , table student terisi otomatis dibagian user_id
}

type StudentAdvisorRequest struct {
	AdvisorId string `json:"advisor_id"`
}
