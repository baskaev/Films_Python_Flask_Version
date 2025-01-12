import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import MoviesList from './components/MoviesList';
import MoviePlayer from './components/MoviePlayer';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MoviesList />} />
        <Route path="/movie/:code" element={<MoviePlayer />} />
      </Routes>
    </Router>
  );
};

export default App;
