#!/bin/bash

# Connection details
DB_HOST="localhost"
DB_USER="user"
DB_PASSWORD="password"
DB_NAME="films_db"

# Export password for non-interactive mode
export PGPASSWORD=$DB_PASSWORD

# Create table SQL for movies
CREATE_MOVIES_TABLE_SQL="""
CREATE TABLE IF NOT EXISTS movies (
    code VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    rating VARCHAR(10),
    year VARCHAR(4),
    image_link TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
"""

# Insert data SQL for movies
INSERT_MOVIES_DATA_SQL="""
INSERT INTO movies (code, title, rating, year, image_link) VALUES
('tt0068646', 'The Godfather', '9.2', '1972', 'https://m.media-amazon.com/images/I/51GrHPaq8QL._AC_.jpg'),
('tt0468569', 'The Dark Knight', '9.0', '2008', 'https://m.media-amazon.com/images/I/71gRxqSToZL._AC_SY679_.jpg'),
('tt0071562', 'The Godfather: Part II', '9.0', '1974', 'https://m.media-amazon.com/images/I/51XItZ7Z7mL._AC_.jpg'),
('tt0120338', 'The Lord of the Rings: The Fellowship of the Ring', '8.8', '2001', 'https://m.media-amazon.com/images/I/71Zb97D4TAL._AC_SY679_.jpg'),
('tt0137523', 'Fight Club', '5', '1999', 'https://m.media-amazon.com/images/I/81kq9Nd9IhL._AC_SY679_.jpg'),
('tt0108052', 'Schindler''s List', '4', '1993', 'https://m.media-amazon.com/images/I/51rJSkRHTRL._AC_.jpg'),
('tt1396484', 'Avengers: Infinity War', '3', '2018', 'https://m.media-amazon.com/images/I/91bQY+opxyL._AC_SY679_.jpg'),
('tt1375666', 'Inception', '2', '2010', 'https://m.media-amazon.com/images/I/71V5GzlfRSL._AC_SY679_.jpg'),
('tt0110912', 'Pulp Fiction', '1', '1994', 'https://m.media-amazon.com/images/I/91zZdeK9GJL._AC_SY679_.jpg')
ON CONFLICT (code) DO NOTHING;
"""

# Create table SQL for tasks
CREATE_TASKS_TABLE_SQL="""
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    isTimerUsed BOOLEAN DEFAULT FALSE,
    runInTime TIMESTAMPTZ,
    priority INT DEFAULT 0,
    paramsJson JSONB,
    done_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
"""

# Execute commands
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "$CREATE_MOVIES_TABLE_SQL"
if [ $? -ne 0 ]; then
  echo "Failed to create movies table."
  exit 1
fi

psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "$INSERT_MOVIES_DATA_SQL"
if [ $? -ne 0 ]; then
  echo "Failed to insert movies data."
  exit 1
fi

psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "$CREATE_TASKS_TABLE_SQL"
if [ $? -ne 0 ]; then
  echo "Failed to create tasks table."
  exit 1
fi

echo "SQL commands executed successfully."
