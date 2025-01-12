import React, { useEffect, useState } from 'react';
import { fetchTopRatedMovies, Movie } from '../api';
import './MoviesList.css';

const MoviesList: React.FC = () => {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchTopRatedMovies()
      .then((data) => setMovies(data))
      .catch((error) => {
        console.error('Ошибка загрузки фильмов:', error);
        setError('Не удалось загрузить фильмы');
      });
  }, []);

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="movies-container">
      {movies.map((movie) => (
        <a
          href={`/movie/${movie.code}`}
          key={movie.code}
          className="movie-item"
        >
          <img src={movie.image_link} alt={movie.title} />
          <div className="movie-details">
            <h3>{movie.title}</h3>
            <p>Рейтинг: {movie.rating}</p>
            <p>Год: {movie.year}</p>
          </div>
        </a>
      ))}
    </div>
  );
};

export default MoviesList;
