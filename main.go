package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Tugas struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	NamaLengkap         string `json:"nama_lengkap"`
	TanggalTugasSelesai string `json:"tanggal_tugas_selesai"`
	LinkTugas           string `json:"link_tugas"`

	CreateAt string
}

var db *gorm.DB

func initDB() {
	dsn := "root@tcp(127.0.0.1:3306)/tugas?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto migrate
	db.AutoMigrate(&Tugas{})
}

func main() {
	initDB()

	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	log.Println("Server running on port 8080...")
	if err := server.ListenAndServe(":5454"); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/tugas":
		switch string(ctx.Method()) {
		case "GET":
			getStudents(ctx)
		case "POST":
			CreateTugas(ctx)
		}
	case "/students/{id}":
		id := ctx.UserValue("id").(string)
		switch string(ctx.Method()) {
		case "GET":
			getStudent(ctx, id)
		case "DELETE":
			deleteStudent(ctx, id)
		}
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func getStudents(ctx *fasthttp.RequestCtx) {
	var students []Tugas
	if err := db.Find(&students).Error; err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	jsonResponse(ctx, students)
}

func CreateTugas(ctx *fasthttp.RequestCtx) {
	var Tugas Tugas
	if err := json.Unmarshal(ctx.PostBody(), &Tugas); err != nil {
		ctx.Error("Invalid input", fasthttp.StatusBadRequest)
		return
	}

	if err := db.Create(&Tugas).Error; err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	jsonResponse(ctx, Tugas)
}

func getStudent(ctx *fasthttp.RequestCtx, id string) {
	var student Tugas
	if err := db.First(&student, id).Error; err != nil {
		ctx.Error("Student not found", fasthttp.StatusNotFound)
		return
	}

	jsonResponse(ctx, student)
}

func deleteStudent(ctx *fasthttp.RequestCtx, id string) {
	if err := db.Delete(&Tugas{}, id).Error; err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func jsonResponse(ctx *fasthttp.RequestCtx, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		ctx.Error("Failed to marshal JSON", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}
