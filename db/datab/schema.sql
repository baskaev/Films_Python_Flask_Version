CREATE TABLE movies (
    code VARCHAR(255) PRIMARY KEY, 
    title VARCHAR(255) NOT NULL,   
    rating VARCHAR(10),            
    year VARCHAR(4),               
    image_link TEXT,                
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,            -- Уникальный идентификатор задачи
    task_name VARCHAR(255) NOT NULL,  -- Название задачи
    isTimerUsed BOOLEAN DEFAULT FALSE,-- Флаг использования таймера
    runInTime TIMESTAMPTZ,            -- Время выполнения задачи
    priority INT DEFAULT 0,           -- Приоритет задачи (0 по умолчанию)
    paramsJson JSONB,                 -- Дополнительные параметры в формате JSON
    done_at TIMESTAMPTZ,               -- Время завершения задачи
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP -- Время создания
);




INSERT INTO movies (code, title, rating, year, image_link) VALUES
('tt0068646', 'The Godfather', '9.2', '1972', 'https://m.media-amazon.com/images/I/51GrHPaq8QL._AC_.jpg'),
('tt0468569', 'The Dark Knight', '9.0', '2008', 'https://m.media-amazon.com/images/I/71gRxqSToZL._AC_SY679_.jpg'),
('tt0071562', 'The Godfather: Part II', '9.0', '1974', 'https://m.media-amazon.com/images/I/51XItZ7Z7mL._AC_.jpg'),
('tt0120338', 'The Lord of the Rings: The Fellowship of the Ring', '8.8', '2001', 'https://m.media-amazon.com/images/I/71Zb97D4TAL._AC_SY679_.jpg'),
('tt0137523', 'Fight Club', '5', '1999', 'https://m.media-amazon.com/images/I/81kq9Nd9IhL._AC_SY679_.jpg'),
('tt0108052', 'Schindler''s List', '4', '1993', 'https://m.media-amazon.com/images/I/51rJSkRHTRL._AC_.jpg'),
('tt1396484', 'Avengers: Infinity War', '3', '2018', 'https://m.media-amazon.com/images/I/91bQY+opxyL._AC_SY679_.jpg'),
('tt1375666', 'Inception', '2', '2010', 'https://m.media-amazon.com/images/I/71V5GzlfRSL._AC_SY679_.jpg'),
('tt0110912', 'Pulp Fiction', '1', '1994', 'https://m.media-amazon.com/images/I/91zZdeK9GJL._AC_SY679_.jpg');

