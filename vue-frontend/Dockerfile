# Используем официальный образ Node.js для сборки
FROM node:18-alpine as build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем весь исходный код и собираем проект
COPY . .
#RUN npm run serve
RUN npm run build

# Используем Nginx для обслуживания статических файлов
FROM nginx:alpine


# Копируем собранные файлы в директорию Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем файл конфигурации Nginx
#COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf
# Открываем порт для Nginx
EXPOSE 80

# Nginx запускается по умолчанию
