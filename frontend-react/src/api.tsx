import axios from 'axios';

const API_URL = 'http://localhost:8082/api';

export interface Movie {
  code: string;
  title: string;
  image_link: string;
  rating: string;
  year: string;
}

export const fetchTopRatedMovies = async (): Promise<Movie[]> => {
  const response = await axios.get<Movie[]>(`${API_URL}/FetchLatestTopRatedMovies`);
  return response.data;
};
