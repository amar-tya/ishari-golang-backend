package entity

import "time"

// RefreshToken represents a stored refresh token for a user session
type RefreshToken struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	TokenHash string     `json:"-" gorm:"not null;index"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
	IPAddress string     `json:"ip_address,omitempty"`

	// Relation
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// IsRevoked checks if the token has been revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsExpired checks if the token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
