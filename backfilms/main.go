package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/baskaev/FILMS_GO_REACT_DOCKER/db/datab"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Инициализация базы данных
	if err := datab.InitDB(); err != nil {
		panic(err)
	}

	// Создаем экземпляр Echo
	e := echo.New()

	// Используем CORS middleware
	e.Use(middleware.CORS())

	// Определяем маршруты
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// Обработчик для маршрута /api
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "API is working!",
		})
	})

	e.GET("/api/FetchMovies", func(c echo.Context) error {
		movies, err := datab.FetchMovies()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, movies)
	})

	e.POST("/api/AddMovie", func(c echo.Context) error {
		code := c.QueryParam("code")
		title := c.QueryParam("title")
		rating := c.QueryParam("rating")
		year := c.QueryParam("year")
		imageLink := c.QueryParam("image_link")

		if code == "" || title == "" || rating == "" || year == "" || imageLink == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing required parameters"})
		}

		movie := datab.Movie{
			Code:      code,
			Title:     title,
			Rating:    rating,
			Year:      year,
			ImageLink: imageLink,
		}

		if err := datab.AddMovie(movie); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Movie added successfully!"})
	})

	e.GET("/api/FetchLatestTopRatedMovies", func(c echo.Context) error {
		movies, err := datab.FetchLatestTopRatedMovies()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, movies)
	})

	e.GET("/api/GetMovieByCode", func(c echo.Context) error {
		code := c.QueryParam("code")
		if code == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing required parameter: code"})
		}

		movie, err := datab.GetMovieByCode(code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, movie)
	})

	e.GET("/api/SearchMovies", func(c echo.Context) error {
		query := c.QueryParam("query")
		yearParam := c.QueryParam("year")
		minRating := c.QueryParam("minRating")

		// Преобразуем параметр `year` в массив
		var years []string
		if yearParam != "" {
			years = strings.Split(yearParam, ",") // Разделяем строку на массив годов
		}

		// Парсим `minRating` в float
		var minRatingFloat float64
		if minRating != "" {
			var err error
			minRatingFloat, err = strconv.ParseFloat(minRating, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid minRating"})
			}
		}

		// Вызов функции поиска фильмов
		movies, err := datab.SearchMovies(query, years, minRatingFloat)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Возвращаем результат
		return c.JSON(http.StatusOK, movies)
	})

	e.GET("/api/AddTaskImdbParser", func(c echo.Context) error {
		fmt.Println("AddTaskImdbParser called")
		task := datab.Task{
			TaskName:    "imdb_parser",
			IsTimerUsed: true,
			RunInTime:   sql.NullTime{},
			Priority:    1,
			ParamsJson:  `{"query": "The Movie", "years": ["2000", "2005", "2010"], "minRating": 6.5}`,
			DoneAt:      sql.NullTime{},
		}
		newID, err := datab.AddTask(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Task added successfully!", "task_id": newID})
	})

	// напиши обработчик для /api/AddMovie?code=tt123456&title=The+Movie&rating=5.0&year=2021&image_link=http://example.com/image.jpg
	// внутри создается объект  datab.Movie и вызывается AddMovie, если все норм результат ок, иначе ошибка

	//  напиши обработчик для /api/FetchLatestTopRatedMovies
	//  возвращает json с переданными фильмами

	//напиши обработчик для /api/GetMovieByCode?code=tt123456
	// возвращает json с переданным фильмом

	// напиши обработчик для /api/SearchMovies?query=The+Movie&years=2000,2005,2010&minRating=6.5

	// Запускаем сервер
	e.Logger.Fatal(e.Start(":8081"))
}
