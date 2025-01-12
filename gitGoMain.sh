#!/bin/bash

# Папка, в которой находятся ваши проекты
base_dir="$HOME/Рабочий стол/filmsProject"

# Перебираем все подпапки в base_dir
for dir in "$base_dir"/*; do
  # Проверяем, что это директория
  if [ -d "$dir" ]; then
    echo "Заходим в директорию: $dir"
    cd "$dir" || continue
    # Проверяем, является ли директория репозиторием git
    if [ -d ".git" ]; then
      echo "Выполняем git checkout main в $dir"
      git checkout main
    else
      echo "Не git репозиторий в $dir"
    fi
    cd "$base_dir" || continue
  fi
done
