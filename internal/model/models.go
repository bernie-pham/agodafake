package model

import (
	"time"

	"github.com/bernie-pham/agodafake/internal/forms"
)

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Phone       string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	RoomName  string
	PriceID   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RoomRestriction struct {
	ID            int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	RestrictionID int
	ReservationID int
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

type Price struct {
	ID         int
	PriceValue int
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
