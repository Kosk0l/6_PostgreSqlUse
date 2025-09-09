package main

import (
	"log"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres" 
	"gorm.io/gorm"            
)


// Структура - Запрос
type Message struct {
	ID 		int 	`json:"id"`
	Text 	string 	`json:"text"`
}

// Структура - Ответ // Что вернет сервер клиенту
type Responce struct {
	Status 	string 	`json:"status"`
	Message string 	`json:"message"`
}

// Переменная базы данных
var db *gorm.DB

func InitDB() {
	// Стринговая строка подключения к базе данных
	dsn := "host=localhost user=postgres password=password dbname=postgres port=port sslmode=disable"
	var err error

	// Принимает информацию о Нужной БД и сама бд принимает параметры подключения // Конфигурации nil
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { //							 // Передача значения в строку
		log.Fatal("Не удалось подключиться к БД: %v", err)
	}

	// Передача структуры в бд для формирования Таблицы по структуре
	db.AutoMigrate(&Message{})
}

//==============================================================================================================

// Просмотр всех имеющихся JSON Структур
func GetJSON(c echo.Context) error { 
	var messages []Message

	// SELECT * FROM messages // Запишет много объектов в массив структур из БД
	if err := db.Find(&messages).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not Find(GET) the Messages",
		})
	}

	// Статус-Код и структуру
	return c.JSON(http.StatusOK, &messages) // Клиент получит JSON-объекты из БД
}

// Создание JSON структуры по данным, полученным фронтом.
func PostJSON(c echo.Context) error {
	var mes Message // Объект класса Message

	// Проверка на ошибку в c.Bind()
	if err := c.Bind(&mes); err != nil { // Если ошибка есть, ТО ..
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not add the Message",
		})
	}

	if err := db.Create(&mes).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not Create the Message in DB",
		})
	}

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: mes.Text,
	})
}

// Обновление данных
func PatchJSON(c echo.Context) error {
	idStr := c.Param("id") // Параметр из URL // приходит в стринге
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Bad Parameter URL",
		})
	}

	// Заполняем параметры стурктуры объекта
	var UpdateMes Message
	if err := c.Bind(&UpdateMes); err != nil { // Если ошибка есть, ТО ..
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Invalid input",
		})
	}

	// Передаем модель структуры.  // Ищем его по id   // Обновляем поле text на UPDT.text // Передаем ошибку
	if err := db.Model(&Message{}).Where("id = ?", id).Update("text", UpdateMes.Text).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Could not UPDATE the MESSAGE",
		})
	}

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: "Message was updated",
	})
}

func DeleteJSON(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Bad Parameter URL",
		})
	}

	if err := db.Delete(&Message{}, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Responce{
			Status: "Error",
			Message: "Invalid id",
		})
	}

	return c.JSON(http.StatusOK, Responce{
		Status: "Success",
		Message: "Message was Deleted",
	})

}

func main() {
	InitDB()
	e := echo.New()

	e.GET("/get", GetJSON)
	e.POST("/post", PostJSON)
	e.PATCH("/patch/:id", PatchJSON)
	e.DELETE("/delete/:id", DeleteJSON)

	e.Start(":8080")
}
