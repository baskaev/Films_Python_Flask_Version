#!/bin/bash

repos=(
  "https://github.com/baskaev/nginx"
  "https://github.com/baskaev/vue-frontend"
  "https://github.com/baskaev/db"
  "https://github.com/baskaev/backfilms"
)

# Папка для хранения клонированных репозиториев
dest_folder="$HOME/Рабочий стол/filmsProject"

# Переходим в папку, где будут клонироваться репозитории
cd "$dest_folder"

# Клонирование каждого репозитория из списка
for repo in "${repos[@]}"
do
  git clone "$repo"
done

echo "Клонирование завершено!"
