package datab

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

type Movie struct {
	Title     string
	Code      string
	Rating    string
	Year      string
	ImageLink string
}

type Task struct {
	ID          int
	TaskName    string
	IsTimerUsed bool
	RunInTime   sql.NullTime
	Priority    int
	ParamsJson  string
	CreatedAt   string
	DoneAt      sql.NullTime
}

// InitDB initializes the database connection
func InitDB() error {
	var err error
	connStr := "host=db user=user password=password dbname=films_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка подключения
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

// FetchMovies retrieves all movies from the database
func FetchMovies() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT code, title, rating, year, image_link FROM movies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var code, title, rating, year, imageLink string
		if err := rows.Scan(&code, &title, &rating, &year, &imageLink); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		movies = append(movies, map[string]interface{}{
			"code":       code,
			"title":      title,
			"rating":     rating,
			"year":       year,
			"image_link": imageLink,
		})
	}
	return movies, nil
}

// AddMovie inserts a new movie into the database
func AddMovie(movie Movie) error {
	// SQL query to insert a new movie
	query := `INSERT INTO movies (code, title, rating, year, image_link) 
			  VALUES ($1, $2, $3, $4, $5)`

	// Execute the query with the movie data
	_, err := db.Exec(query, movie.Code, movie.Title, movie.Rating, movie.Year, movie.ImageLink)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	return nil
}

func AddTask(task Task) (int, error) {
	query := `
        INSERT INTO tasks (task_name, isTimerUsed, runInTime, priority, paramsJson, done_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	var newID int
	err := db.QueryRow(query,
		task.TaskName,
		task.IsTimerUsed,
		task.RunInTime,
		task.Priority,
		task.ParamsJson,
		task.DoneAt,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}
	return newID, nil
}

// DeleteTaskByID deletes a task from the database by its ID.
func DeleteTaskByID(taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	// Execute the query to delete the task
	result, err := db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task with id %d: %w", taskID, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %d", taskID)
	}

	return nil
}

// FetchLatestTopRatedMovies retrieves the last 50 (or fewer) movies added with IMDb rating greater than 6
func FetchLatestTopRatedMovies() ([]map[string]interface{}, error) {
	query := `
		SELECT code, title, rating, year, image_link
		FROM movies
		WHERE rating::float > 6
		ORDER BY created_at DESC, rating::float DESC
		LIMIT 50;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest top-rated movies: %w", err)
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var code, title, rating, year, imageLink string
		if err := rows.Scan(&code, &title, &rating, &year, &imageLink); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		movies = append(movies, map[string]interface{}{
			"code":       code,
			"title":      title,
			"rating":     rating,
			"year":       year,
			"image_link": imageLink,
		})
	}

	// Проверяем, произошли ли ошибки при итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return movies, nil
}

// GetMovieByCode retrieves a movie from the database by its unique code.
func GetMovieByCode(code string) (*Movie, error) {
	query := `SELECT code, title, rating, year, image_link FROM movies WHERE code = $1`
	var movie Movie
	err := db.QueryRow(query, code).Scan(&movie.Code, &movie.Title, &movie.Rating, &movie.Year, &movie.ImageLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to fetch movie: %w", err)
	}
	return &movie, nil
}

func SearchMovies(query string, years []string, minRating float64) ([]Movie, error) {
	var movies []Movie

	// Формируем SQL-запрос
	sqlQuery := `SELECT code, title, rating, year, image_link FROM movies WHERE title ILIKE $1`
	var args []interface{}
	args = append(args, "%"+query+"%")

	if len(years) > 0 {
		sqlQuery += " AND year = ANY($2)"    // Используем ANY для работы с массивами
		args = append(args, pq.Array(years)) // Преобразуем срез в массив для PostgreSQL
	}

	if minRating > 0 {
		sqlQuery += " AND rating >= $3"
		args = append(args, minRating)
	}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Code, &movie.Title, &movie.Rating, &movie.Year, &movie.ImageLink); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

// FetchTopPriorityTask retrieves the task with isTimerUsed=true, highest priority, and earliest creation time.
func FetchTopPriorityTask() (Task, error) {
	var task Task

	query := `
        SELECT id, task_name, isTimerUsed, runInTime, priority, paramsJson, created_at, done_at
        FROM tasks
        WHERE isTimerUsed = true
        ORDER BY priority DESC, created_at ASC
        LIMIT 1;
    `

	row := db.QueryRow(query)

	err := row.Scan(&task.ID, &task.TaskName, &task.IsTimerUsed, &task.RunInTime, &task.Priority, &task.ParamsJson, &task.CreatedAt, &task.DoneAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// No task found that matches the criteria
			return Task{}, fmt.Errorf("no task found with isTimerUsed=true")
		}
		return Task{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return task, nil
}

// FetchAllTasks retrieves all tasks from the database.
func FetchAllTasks() ([]Task, error) {
	query := "SELECT id, task_name, isTimerUsed, runInTime, priority, paramsJson, created_at, done_at FROM tasks;"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.TaskName, &task.IsTimerUsed, &task.RunInTime, &task.Priority, &task.ParamsJson, &task.CreatedAt, &task.DoneAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}
