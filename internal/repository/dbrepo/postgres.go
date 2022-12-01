package dbrepo

import (
	"context"
	"time"

	"github.com/bernie-pham/agodafake/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}
func (m *postgresDBRepo) InsertReservation(res model.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rsvID int
	stmt := `insert into "Reservations" (first_name, last_name, email,
		 phone, start_date, end_date, room_id)
		 values ($1, $2, $3, $4, $5, $6, $7) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
	).Scan(&rsvID)

	if err != nil {
		return -1, err
	}

	if err != nil {
		return -1, err
	}
	return int(rsvID), nil
}

func (m *postgresDBRepo) InsertRoomRestriction(rr model.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into "Room_restriction" (start_date, end_date, room_id,
		reservation_id, restriction_id)
		 values ($1, $2, $3, $4, $5)`
	_, err := m.DB.ExecContext(ctx, stmt,
		rr.StartDate,
		rr.EndDate,
		rr.RoomID,
		rr.ReservationID,
		rr.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) SearchRoomAvailableByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numOccupied int

	stmt := `select count(rr.id) from "Room_restriction" rr  
		where room_id = $1 and
		$2 < rr.end_date and $3 > rr.start_date`

	err := m.DB.QueryRowContext(ctx, stmt, roomID, start, end).Scan(&numOccupied)

	if err != nil {
		return false, err
	}
	if numOccupied > 0 {
		return false, nil
	}
	return true, nil
}
func (m *postgresDBRepo) GetRoomByRoomID(roomID int) (model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room model.Room

	stmt := `select room_name, price_id from "Rooms"  
		where id = $1`

	err := m.DB.QueryRowContext(ctx, stmt, roomID).Scan(&room.RoomName, &room.PriceID)

	if err != nil {
		return room, err
	}

	return room, nil
}

func (m *postgresDBRepo) SearchRoomsAvailableByDates(start, end time.Time) ([]model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.room_name
		from "Rooms" r 
		where r.id not in 
		(select rr.room_id  from "Room_restriction" rr 
		where $1 < rr.end_date and $2 > rr.start_date ) `

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	var rooms []model.Room
	var roomName string
	var roomID int
	for rows.Next() {
		err := rows.Scan(&roomID, &roomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, model.Room{
			RoomName: roomName,
			ID:       roomID,
		})
	}
	return rooms, nil
}

func (m *postgresDBRepo) InsertUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		insert into "Users" (first_name, last_name, email, phone, access_level, password)
		values ($1, $2, $3, $4, $5, $6)
	`

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = m.DB.ExecContext(
		ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.AccessLevel,
		hashedpassword,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetUserByID(userID int) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user model.User

	query := `
		select id, first_name, last_name, email, phone, created_at, updated_at, access_level
		from "Users"
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.AccessLevel,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (m *postgresDBRepo) UpdateUserByID(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update "Users" 
		set first_name = $1, last_name = $2, access_level = $3, phone = $4, password = $5
	`

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(
		ctx, query,
		user.FirstName,
		user.LastName,
		user.AccessLevel,
		user.Phone,
		hashedpassword,
	)
	if err != nil {
		return err
	}
	return nil
}
func (m *postgresDBRepo) VerifyUser(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var storedPassword string
	var userID int
	query := `
		select password, id 
		from "Users"
		where email = $1
	`
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&storedPassword, &userID)

	if err != nil {
		return 0, "", err
	}

	// Then moving to authenticate step
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", err
	} else if err != nil {
		return 0, "", err
	}

	return userID, "", nil
}
