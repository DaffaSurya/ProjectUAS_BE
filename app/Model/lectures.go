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
