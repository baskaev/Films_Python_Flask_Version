from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import func, and_, or_
from datetime import datetime
import json

db = SQLAlchemy()

class Movie(db.Model):
    __tablename__ = 'movies'
    code = db.Column(db.String(255), primary_key=True)  # Первичный ключ
    title = db.Column(db.String(255), nullable=False)
    rating = db.Column(db.String(10))
    year = db.Column(db.String(4))
    image_link = db.Column(db.Text)
    created_at = db.Column(db.DateTime, default=db.func.current_timestamp())

class Task(db.Model):
    __tablename__ = 'tasks'
    id = db.Column(db.Integer, primary_key=True)
    task_name = db.Column(db.String(255), nullable=False)
    is_timer_used = db.Column(db.Boolean, name='istimerused', default=False)  # Исправлено имя столбца
    run_in_time = db.Column(db.String, name='runintime')  # Исправлено имя и тип столбца
    priority = db.Column(db.Integer, default=0)
    params_json = db.Column(db.String, name='paramsjson')  # Исправлено имя и тип столбца
    done_at = db.Column(db.DateTime)
    created_at = db.Column(db.DateTime, default=db.func.current_timestamp())

def init_db(app):
    """Инициализация базы данных."""
    db.init_app(app)
    with app.app_context():
        db.create_all()

def fetch_movies():
    """Получить все фильмы."""
    movies = Movie.query.all()
    return [{
        "code": movie.code,
        "title": movie.title,
        "rating": movie.rating,
        "year": movie.year,
        "image_link": movie.image_link
    } for movie in movies]



def add_movie(movie_data):
    """Добавить фильм."""
    new_movie = Movie(
        code=movie_data['code'],
        title=movie_data['title'],
        rating=movie_data['rating'],
        year=movie_data['year'],
        image_link=movie_data['image_link']
    )
    db.session.add(new_movie)
    db.session.commit()

def fetch_latest_top_rated_movies():
    """Получить последние 50 фильмов с рейтингом выше 6."""
    movies = Movie.query.filter(Movie.rating.cast(db.Float) > 6.0) \
                        .order_by(Movie.created_at.desc(), Movie.rating.desc()) \
                        .limit(50).all()
    return [{
        "code": movie.code,
        "title": movie.title,
        "rating": movie.rating,
        "year": movie.year,
        "image_link": movie.image_link
    } for movie in movies]

def get_movie_by_code(code):
    """Получить фильм по коду."""
    movie = Movie.query.filter_by(code=code).first()
    if movie:
        return {
            "code": movie.code,
            "title": movie.title,
            "rating": movie.rating,
            "year": movie.year,
            "image_link": movie.image_link
        }
    return None


def search_movies(query, years, min_rating):
    """Поиск фильмов."""
    query_filter = Movie.title.ilike(f'%{query}%')
    if years:
        query_filter = and_(query_filter, Movie.year.in_(years))
    if min_rating > 0:
        query_filter = and_(query_filter, Movie.rating.cast(db.Float) >= min_rating)

    movies = Movie.query.filter(query_filter).all()
    return [{
        "code": movie.code,
        "title": movie.title,
        "rating": movie.rating,
        "year": movie.year,
        "image_link": movie.image_link
    } for movie in movies]

def add_task(task):
    """Добавить задачу."""
    new_task = Task(
        task_name=task['task_name'],
        is_timer_used=task['is_timer_used'],
        run_in_time=task.get('run_in_time'),
        priority=task['priority'],
        params_json=json.dumps(task['params_json']),  # Преобразуем JSON в строку
        done_at=task.get('done_at')
    )
    db.session.add(new_task)
    db.session.commit()
    return new_task.id

def delete_task_by_id(task_id):
    """Удалить задачу по ID."""
    task = Task.query.get(task_id)
    if task:
        db.session.delete(task)
        db.session.commit()
        return True
    return False

def fetch_top_priority_task():
    """Получить задачу с наивысшим приоритетом."""
    task = Task.query.filter_by(is_timer_used=True) \
                     .order_by(Task.priority.desc(), Task.created_at.asc()) \
                     .first()
    if task:
        return {
            "id": task.id,
            "task_name": task.task_name,
            "is_timer_used": task.is_timer_used,
            "run_in_time": task.run_in_time,
            "priority": task.priority,
            "params_json": json.loads(task.params_json) if task.params_json else None,  # Преобразуем строку в JSON
            "created_at": task.created_at,
            "done_at": task.done_at
        }
    return None

def fetch_all_tasks():
    """Получить все задачи."""
    tasks = Task.query.all()
    return [{
        "id": task.id,
        "task_name": task.task_name,
        "is_timer_used": task.is_timer_used,
        "run_in_time": task.run_in_time,
        "priority": task.priority,
        "params_json": json.loads(task.params_json) if task.params_json else None,  # Преобразуем строку в JSON
        "created_at": task.created_at,
        "done_at": task.done_at
    } for task in tasks]