package repository

import (
	"time"

	"github.com/bernie-pham/agodafake/internal/model"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res model.Reservation) (int, error)
	InsertRoomRestriction(rr model.RoomRestriction) error
	SearchRoomAvailableByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchRoomsAvailableByDates(start, end time.Time) ([]model.Room, error)
	GetRoomByRoomID(roomID int) (model.Room, error)

	GetUserByID(userID int) (model.User, error)
	UpdateUserByID(user model.User) error
	VerifyUser(email, password string) (int, string, error)
	InsertUser(user model.User) error
}
