package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/baskaev/FILMS_GO_REACT_DOCKER/db/datab"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const (
	imdbParser = "imdb_parser"
)

type imdbParserParams struct {
	Size int `json:"quentity"`
}

func main() {
	// Инициализация базы данных
	if err := datab.InitDB(); err != nil {
		panic(err)
	}

	// получить из базы данных все записи в таблице tasks, затем взять только те у которых стоит таймер true, затем взять
	// те у кого выше приоритет, затем если приоритет равный допустим взять ту таску что раньше была добавлена в таблицу

	for {
		// Получение задачи из базы данных
		task, err := datab.FetchTopPriorityTask()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Minute)
			continue
		}

		// Парсинг данных
		switch task.TaskName {
		case imdbParser:
			var params imdbParserParams
			filmCount := 100
			err := json.Unmarshal([]byte(task.ParamsJson), &params)
			if err == nil && params.Size > 0 {
				filmCount = params.Size
			}

			rodRemote := os.Getenv("ROD_REMOTE")
			if rodRemote == "" {
				fmt.Println("ROD_REMOTE is not set")
				break
			}

			l := launcher.MustNewManaged(rodRemote)

			l.Set("disable-gpu").Delete("disable-gpu")

			l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

			browser := rod.New().Client(l.MustClient()).MustConnect()

			launcher.Open(browser.ServeMonitor(""))

			fmt.Println(
				browser.MustPage("https://developer.mozilla.org").MustEval("() => document.title"),
			)

			// Открытие страницы IMDb
			page := browser.MustPage("https://www.imdb.com/search/title/?title_type=feature")

			// Принятие куки
			//page.MustScreenshot("/tmp/debug.png")
			// acceptButton := page.MustElement("button[data-testid='accept-button']") // HERE IS THE PROBLEM!!! when do it in docker container
			// if acceptButton != nil {
			// 	acceptButton.MustClick()
			// 	fmt.Println("Принято использование куки.")
			// 	time.Sleep(1 * time.Second)
			// }

			// Ждем, пока элементы загрузятся
			page.MustWaitLoad()

			// Список для хранения фильмов
			var movies []datab.Movie
			iter := 0
			for {
				// Получаем все элементы фильмов, я пока не понял как игнорировать уже скачанные фильмы
				items := page.MustElements("li.ipc-metadata-list-summary-item")
				if len(items) > filmCount {
					break
				}
				// Проходим по каждому элементу и извлекаем данные
				for i := iter; i < len(items); i++ { //item := range items {
					item := items[i]
					// Название фильма
					title := "Неизвестно"
					has, titleElement, _ := item.Has(".ipc-title__text")
					if has {
						title = removeNumberPrefix(titleElement.MustText())
					}

					// Уникальный код фильма
					code := "Неизвестно"
					has, codeElement, _ := item.Has(".ipc-title-link-wrapper")
					if has {
						codeAttr := codeElement.MustAttribute("href")
						if codeAttr != nil && len(*codeAttr) > 16 {
							code = (*codeAttr)[7:16] // Извлекаем код фильма из ссылки (например, tt1262426)
						}
					}

					// Год фильма
					year := "0"
					has, yearElement, _ := item.Has(".sc-300a8231-7")
					if has {
						year = isNil(yearElement)
					}

					// Рейтинг фильма
					rating := "0"
					has, ratingElement, _ := item.Has(".ipc-rating-star--rating")
					if has {
						rating = isNil(ratingElement)
					}

					// Ссылка на картинку
					imageLink := "Не найдено"
					has, imageElement, _ := item.Has("img.ipc-image")
					if has {
						imageAttr := imageElement.MustAttribute("src")
						if imageAttr != nil {
							imageLink = *imageAttr
						}
					}

					// Создаем объект фильма и добавляем в список
					movie := datab.Movie{
						Title:     title,
						Code:      code,
						Rating:    rating,
						Year:      year,
						ImageLink: imageLink,
					}
					movies = append(movies, movie)
					datab.AddMovie(movie)

					// Выводим информацию о фильме
					fmt.Printf("Название: %s\nКод: %s\nРейтинг: %s\nГод: %s\nСсылка на картинку: %s\n------------\n", movie.Title, movie.Code, movie.Rating, movie.Year, movie.ImageLink)
				}

				iter = len(items)
				// Проверяем наличие кнопки "50 more"
				buttons := page.MustElements("button.ipc-see-more__button")
				if len(buttons) > 0 {
					// Прокрутка до кнопки "50 more"
					buttons[0].MustScrollIntoView()

					// Нажатие на кнопку "50 more"
					buttons[0].MustClick()
					fmt.Println("Нажата кнопка '50 more', загружаются новые фильмы...")
					time.Sleep(5 * time.Second)
				} else {
					fmt.Println("Кнопка '50 more' не найдена. Завершаем парсинг.")
					break
				}
			}

			// Закрытие браузера
			browser.MustClose()

			// Выводим все фильмы из списка
			// fmt.Println("Все фильмы:")
			// for _, movie := range movies {
			// 	fmt.Printf("Название: %s\nКод: %s\nРейтинг: %s\nГод: %s\nСсылка на картинку: %s\n\n", movie.Title, movie.Code, movie.Rating, movie.Year, movie.ImageLink)
			// }

		default:
			fmt.Println("Another task")
		}

		datab.DeleteTaskByID(task.ID)

		time.Sleep(time.Minute)

	}
}

func isNil(a *rod.Element) string {
	var str string
	if a != nil {
		str = a.MustText()
	} else {
		str = "Неизвестен" // Если элемент не найден, ставим значение по умолчанию
	}
	return str
}

// Убирает номер из начала названия
func removeNumberPrefix(title string) string {
	re := regexp.MustCompile(`^\d+\.\s*`)
	return strings.TrimSpace(re.ReplaceAllString(title, ""))
}
