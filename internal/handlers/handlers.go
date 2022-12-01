package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bernie-pham/agodafake/db/driver"
	"github.com/bernie-pham/agodafake/internal/config"
	"github.com/bernie-pham/agodafake/internal/forms"
	"github.com/bernie-pham/agodafake/internal/helpers"
	"github.com/bernie-pham/agodafake/internal/model"
	"github.com/bernie-pham/agodafake/internal/render"
	"github.com/bernie-pham/agodafake/internal/repository"
	"github.com/bernie-pham/agodafake/internal/repository/dbrepo"
	"github.com/go-chi/chi"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (r *Repository) Home(w http.ResponseWriter, req *http.Request) {
	remoteIP := req.RemoteAddr

	r.App.Session.Put(req.Context(), "remote_ip", remoteIP)

	render.Template(w, req, "home.page.html", &model.TemplateData{})
}
func (r *Repository) About(w http.ResponseWriter, req *http.Request) {

	render.Template(w, req, "about.page.html", &model.TemplateData{})
}

func (r *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	res, ok := r.App.Session.Get(req.Context(), "reservation").(model.Reservation)
	if !ok {
		log.Println("Cannot get item from Session")
		r.App.Session.Put(req.Context(), "error", "Cannot get Reservation from Session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	room, err := r.DB.GetRoomByRoomID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.Room.PriceID = room.PriceID
	r.App.Session.Put(req.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, req, "make-reservation.page.html", &model.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (r *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {
	res, ok := r.App.Session.Get(req.Context(), "reservation").(model.Reservation)
	if !ok {
		helpers.ServerError(w, fmt.Errorf("Cannot get reservation from session"))
		return
	}
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(req.PostForm)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, req)
	form.MinLength("last_name", 3, req)
	form.IsEmail("email", req)

	roomID, err := strconv.Atoi(req.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.FirstName = req.Form.Get("first_name")
	res.LastName = req.Form.Get("last_name")
	res.Email = req.Form.Get("email")
	res.Phone = req.Form.Get("phone")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res
		render.Template(w, req, "make-reservation.page.html", &model.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	rsvID, err := r.DB.InsertReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
	}
	roomRestriction := model.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        roomID,
		RestrictionID: 1,
		ReservationID: rsvID,
	}
	log.Println("Reservation ID: ", rsvID)
	err = r.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		helpers.ServerError(w, err)
	}

	mailReq := model.MailData{
		To:      res.Email,
		From:    "",
		Subject: "Make Reservation Successfully",
		Content: fmt.Sprintf(
			"<h3>Hello %s %s,</h3></br>"+
				"<p>You've made successfully the reservation</p></br>"+
				"<p>Your Reservation information:</p></br>"+
				"<ul>"+
				"<li>Room Name: %s</li>"+
				"<li>Arrival: %s</li>"+
				"<li>Departure: %s</li></lu></br>"+
				"<p>Thank you for your trust</p>",
			res.FirstName, res.LastName,
			res.Room.RoomName,
			res.StartDate,
			res.EndDate,
		),
	}
	r.App.MailChan <- mailReq

	r.App.Session.Put(req.Context(), "reservation", res)
	http.Redirect(w, req, "/reservation-summary", http.StatusSeeOther)
}

func (r *Repository) Generals(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "generals.page.html", &model.TemplateData{})
}
func (r *Repository) Majors(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "majors.page.html", &model.TemplateData{})
}
func (r *Repository) Availability(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "search-availability.page.html", &model.TemplateData{})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (r *Repository) AvailabilityJSON(w http.ResponseWriter, req *http.Request) {
	sd := req.Form.Get("start")
	ed := req.Form.Get("end")

	roomID, err := strconv.Atoi(req.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	fmt.Println(startDate, endDate)
	isAvailable, err := r.DB.SearchRoomAvailableByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK:        isAvailable,
		Message:   "!!",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}
	fmt.Println(resp)

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (r *Repository) PostAvailability(w http.ResponseWriter, req *http.Request) {
	start_date := req.Form.Get("start_date")
	end_date := req.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start_date)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end_date)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := r.DB.SearchRoomsAvailableByDates(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		r.App.Session.Put(req.Context(), "error", "No Availability")
		http.Redirect(w, req, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := model.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	r.App.Session.Put(req.Context(), "reservation", reservation)
	render.Template(w, req, "rooms-available.page.html", &model.TemplateData{
		Data: data,
	})
}
func (r *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "contact.page.html", &model.TemplateData{})
}
func (r *Repository) ShowLogin(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "login.page.html", &model.TemplateData{
		Form: forms.New(nil),
	})
}
func (r *Repository) Logout(w http.ResponseWriter, req *http.Request) {
	_ = r.App.Session.Destroy(req.Context())
	_ = r.App.Session.RenewToken(req.Context())

	r.App.Session.Put(req.Context(), "flash", "Logout Successfully")

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
func (r *Repository) VerifyUser(w http.ResponseWriter, req *http.Request) {
	_ = r.App.Session.RenewToken(req.Context())
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(req.Form)
	form.Required("email", "password")
	form.IsEmail("email", req)
	form.MinLength("password", 6, req)

	if !form.Valid() {
		render.Template(w, req, "login.page.html", &model.TemplateData{
			Form: form,
		})
		return
	}
	email := req.Form.Get("email")
	password := req.Form.Get("password")
	userID, _, err := r.DB.VerifyUser(email, password)

	if err == sql.ErrNoRows {
		r.App.Session.Put(req.Context(), "warning", "Invalid Email or Password, Try again!")
		http.Redirect(w, req, "/user/login", http.StatusSeeOther)
		return
	}
	user, err := r.DB.GetUserByID(userID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	r.App.Session.Put(req.Context(), "flash", fmt.Sprintf("Hello %s %s", user.FirstName, user.LastName))
	r.App.Session.Put(req.Context(), "user", user)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
func (r *Repository) ShowRegister(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "register.page.html", &model.TemplateData{
		Form: forms.New(nil),
	})
}
func (r *Repository) RegisterUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(req.PostForm)
	form.Required("first_name", "last_name", "email", "phone", "password", "re_password")
	form.MinLength("first_name", 3, req)
	form.MinLength("last_name", 3, req)
	form.MinLength("password", 6, req)
	form.IsEqual("password", "re_password", req)
	form.IsEmail("email", req)

	user := model.User{
		FirstName:   req.Form.Get("first_name"),
		LastName:    req.Form.Get("last_name"),
		Email:       req.Form.Get("email"),
		Phone:       req.Form.Get("phone"),
		Password:    req.Form.Get("password"),
		AccessLevel: 1,
	}

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user
		render.Template(w, req, "register.page.html", &model.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = r.DB.InsertUser(user)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, req, "/user/login", http.StatusSeeOther)
}

func (r *Repository) ReservationSummary(w http.ResponseWriter, req *http.Request) {
	reservation, ok := r.App.Session.Get(req.Context(), "reservation").(model.Reservation)
	if !ok {
		log.Println("Cannot get item from Session")
		r.App.Session.Put(req.Context(), "error", "Cannot get Reservation from Session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}
	r.App.Session.Remove(req.Context(), "reservation")
	r.App.Session.Put(req.Context(), "flash", "Make Reservation Successfully")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, req, "reservation-summary.page.html", &model.TemplateData{
		Data: data,
	})
}

func (r *Repository) ChooseRoom(w http.ResponseWriter, req *http.Request) {
	// Get room id from request URL eg: http://localhost:8080/choose-room/12, so it will get 12
	roomID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// We need to get the reservation model from session out, and then add the room id into it,
	// then we update the session with edited reservation.

	res, ok := r.App.Session.Get(req.Context(), "reservation").(model.Reservation)
	if !ok {
		log.Println("Cannot get item from Session")
		r.App.Session.Put(req.Context(), "error", "Cannot get Reservation from Session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}
	res.RoomID = roomID
	r.App.Session.Put(req.Context(), "reservation", res)

	http.Redirect(w, req, "/make-reservation", http.StatusSeeOther)
}
func (r *Repository) BookRoom(w http.ResponseWriter, req *http.Request) {
	// id, sd, ed
	id, _ := strconv.Atoi(req.URL.Query().Get("id"))
	sd := req.URL.Query().Get("sd")
	ed := req.URL.Query().Get("ed")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	res := model.Reservation{
		RoomID:    id,
		StartDate: startDate,
		EndDate:   endDate,
	}
	r.App.Session.Put(req.Context(), "reservation", res)
	http.Redirect(w, req, "/make-reservation", http.StatusSeeOther)
}

// func (r *Repository) RoomsAvailable(w http.ResponseWriter, req *http.Request) {
// 	rooms, ok := r.App.Session.Get(req.Context(), "rooms").([]model.Room)
// 	if !ok {
// 		log.Println("Cannot get item from Session")
// 		r.App.Session.Put(req.Context(), "error", "Cannot get Rooms from Session")
// 		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
// 		return
// 	}
// 	r.App.Session.Remove(req.Context(), "rooms")
// 	r.App.Session.Put(req.Context(), "flash", "Search Rooms Availability Successfully")
// 	data := make(map[string]interface{})
// 	data["rooms"] = rooms
// 	render.Template(w, req, "rooms_available.page.tml", &model.TemplateData{
// 		Data: data,
// 	})
// }

func (r *Repository) AdminDashboard(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "admin-dashboard.page.html", &model.TemplateData{
		Form: forms.New(nil),
	})
}
