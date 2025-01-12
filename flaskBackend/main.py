from flask import Flask, jsonify, request
from flask_cors import CORS  # Импортируем CORS
from datab import db, init_db, fetch_movies, add_movie, fetch_latest_top_rated_movies, \
                  get_movie_by_code, search_movies, add_task, delete_task_by_id, \
                  fetch_top_priority_task, fetch_all_tasks

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'postgresql://user:password@db:5432/films_db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False


CORS(app)


# Инициализация базы данных
init_db(app)

@app.route('/')
def hello_world():
    return 'Hello, Flask!'

@app.route('/api')
def api_status():
    return jsonify({"message": "API is working!"})

@app.route('/api/FetchMovies', methods=['GET'])
def fetch_movies_route():
    movies = fetch_movies()
    return jsonify(movies)


@app.route('/api/AddMovie', methods=['POST'])
def add_movie_route():
    # Получаем данные из query-параметров
    code = request.args.get('code')
    title = request.args.get('title')
    rating = request.args.get('rating')
    year = request.args.get('year')
    image_link = request.args.get('image_link')

    # Проверяем, что все параметры присутствуют
    if not all([code, title, rating, year, image_link]):
        return jsonify({"error": "missing required parameters"}), 400

    # Создаем словарь с данными фильма
    movie_data = {
        'code': code,
        'title': title,
        'rating': rating,
        'year': year,
        'image_link': image_link
    }

    # Добавляем фильм в базу данных
    add_movie(movie_data)
    return jsonify({"message": "Movie added successfully!"}), 200

@app.route('/api/FetchLatestTopRatedMovies', methods=['GET'])
def fetch_latest_top_rated_movies_route():
    movies = fetch_latest_top_rated_movies()
    return jsonify(movies)

@app.route('/api/GetMovieByCode', methods=['GET'])
def get_movie_by_code_route():
    code = request.args.get('code')
    if not code:
        return jsonify({"error": "missing required parameter: code"}), 400

    movie = get_movie_by_code(code)
    if movie:
        return jsonify(movie)
    return jsonify({"error": "Movie not found"}), 404

@app.route('/api/SearchMovies', methods=['GET'])
def search_movies_route():
    # Получаем параметры из запроса
    query = request.args.get('query')
    year_param = request.args.get('year', '')
    min_rating = request.args.get('minRating', '')

    # Преобразуем параметр `year` в список
    years = year_param.split(',') if year_param else []

    # Парсим `minRating` в float
    min_rating_float = 0.0
    if min_rating:
        try:
            min_rating_float = float(min_rating)
        except ValueError:
            return jsonify({"error": "invalid minRating"}), 400

    # Вызов функции поиска фильмов
    try:
        movies = search_movies(query, years, min_rating_float)
    except Exception as e:
        return jsonify({"error": str(e)}), 500

    # Возвращаем результат
    return jsonify(movies)

@app.route('/api/AddTaskImdbParser', methods=['GET'])
def add_task_imdb_parser_route():
    task = {
        "task_name": "imdb_parser",
        "is_timer_used": True,
        "priority": 1,
        "params_json": {"query": "The Movie", "years": ["2000", "2005", "2010"], "minRating": 6.5}
    }
    task_id = add_task(task)
    return jsonify({"message": "Task added successfully!", "task_id": task_id})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)