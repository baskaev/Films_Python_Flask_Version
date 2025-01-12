import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';

const MoviePlayer: React.FC = () => {
  const { code } = useParams<{ code: string }>();

  useEffect(() => {
    const script = document.createElement('script');
    script.src = 'https://kinobox.tv/kinobox.min.js';
    script.async = true;
    script.onload = () => {
      if (window.kbox) {
        window.kbox('.kinobox_player', {
          search: { imdb: code },
          // players: { Cdnmovies: { enable: true }, Kodik: { enable: true } },
          players: {
                
            Cdnmovies: {
                enable: true
            },
            Videocdn: {
                enable: true
            },
            Ashdi: {
                enable: true
            },  
            Collaps: {
                enable: true
            },
            Hdvb: {
                enable: true
            },
            Kodik: {
                enable: true
            },
            Vibix: {
                enable: true
            },
            Voidboost: {
                enable: true
            },
            Turbo: {
                enable: true
            },
            Alloha: {
                enable: true
            },
        },
        });
      }
    };
    document.body.appendChild(script);
  }, [code]);

  return (
    <div>
      <h1>Просмотр фильма</h1>
      <div className="kinobox_player"></div>
    </div>
  );
};

export default MoviePlayer;
