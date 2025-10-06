package models

import "time"

type User struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"` // ID auto-generated, no validation
	SocialID       string `gorm:"uniqueIndex:social_provider" json:"social_id" validate:"required"`
	SocialProvider string `gorm:"uniqueIndex:social_provider" json:"social_provider" validate:"required,oneof=facebook instagram"`
	Name           string `json:"name" validate:"required,min=1,max=100"`
	Email          string `json:"email" validate:"required,email"`
	AvatarURL      string `json:"avatar_url" validate:"omitempty,url"`
	IsAdmin        bool   `json:"is_admin" validate:"-"`
	CreatedAt      int64  `gorm:"type:bigint;autoCreateTime:milli" json:"created_at" validate:"-"`
	UpdatedAt      int64  `gorm:"type:bigint;autoCreateTime:milli" json:"updated_at" validate:"-"`
}

type Admin struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	Username string `gorm:"unique" json:"username" validate:"required,min=3,max=50"`
	Email    string `gorm:"unique" json:"email" validate:"required,email"`
}

type Activity struct {
	ID          int64    `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	Title       string   `json:"title" validate:"required,min=1,max=200"`
	TargetCount int      `json:"target_count" validate:"required,min=1,max=100"`
	LocationID  int64    `json:"location_id" validate:"required,min=1"`
	Description string   `json:"description" validate:"omitempty,max=1000"`
	CreatedBy   int64    `json:"created_by" validate:"required,min=1"`
	Location    Location `gorm:"foreignKey:LocationID" json:"location" validate:"-"`
}

type Location struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Address   string  `json:"address" validate:"required,min=1,max=200"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type Match struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	ActivityID  int64     `json:"activity_id" validate:"required,min=1"`
	OrganizerID int64     `json:"organizer_id" validate:"required,min=1"`
	MatchTime   time.Time `json:"match_time" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=open completed cancelled"`
	Activity    Activity  `gorm:"foreignKey:ActivityID" json:"activity" validate:"-"`
	Organizer   User      `gorm:"foreignKey:OrganizerID" json:"organizer" validate:"-"`
}

type MatchParticipant struct {
	ID       int64     `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	MatchID  int64     `json:"match_id" validate:"required,min=1"`
	UserID   int64     `json:"user_id" validate:"required,min=1"`
	Status   string    `json:"status" validate:"required,oneof=pending approved rejected"`
	JoinedAt time.Time `json:"joined_at" validate:"required"`
	Match    Match     `gorm:"foreignKey:MatchID" json:"match" validate:"-"`
	User     User      `gorm:"foreignKey:UserID" json:"user" validate:"-"`
}

type Review struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	MatchID    int64     `json:"match_id" validate:"required,min=1"`
	ReviewerID int64     `json:"reviewer_id" validate:"required,min=1"`
	RevieweeID int64     `json:"reviewee_id" validate:"required,min=1"`
	Score      int       `json:"score" validate:"required,min=3,max=5"`
	Comment    string    `json:"comment" validate:"omitempty,max=500"`
	CreatedAt  time.Time `json:"created_at" validate:"-"`
	Match      Match     `gorm:"foreignKey:MatchID" json:"match" validate:"-"`
	Reviewer   User      `gorm:"foreignKey:ReviewerID" json:"reviewer" validate:"-"`
	Reviewee   User      `gorm:"foreignKey:RevieweeID" json:"reviewee" validate:"-"`
}

type ReviewLike struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id" validate:"-"`
	ReviewID int64  `json:"review_id" validate:"required,min=1"`
	UserID   int64  `json:"user_id" validate:"required,min=1"`
	IsLike   bool   `json:"is_like" validate:"required"`
	Review   Review `gorm:"foreignKey:ReviewID" json:"review" validate:"-"`
	User     User   `gorm:"foreignKey:UserID" json:"user" validate:"-"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id" validate:"-"`
	UserID    uint      `gorm:"index" json:"user_id" validate:"required"`
	Token     string    `gorm:"not null" json:"-" validate:"required,min=60"` // hashed token
	ExpiresAt time.Time `gorm:"index" json:"expires_at" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user" validate:"-"`
}
