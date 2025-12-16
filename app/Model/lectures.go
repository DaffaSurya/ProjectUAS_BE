package model

// Verify atau reject yang bisa dilakukan oleh dosen wali
type VerifyRequest struct {
	Status          string  `json:"status"` // verified | rejected
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

// verify atau reject yang bisa dilakukan oleh dosen wali

// Verify atau reject yang bisa dilakukan oleh dosen wali
type RejectRequest struct {
	Reason string `json:"reason"`
}

// verify atau reject yang bisa dilakukan oleh dosen wali

type LecturerResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Department string `json:"department"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
}

type AdviseeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
}
