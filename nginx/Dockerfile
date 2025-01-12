
# Используем официальный образ Nginx
FROM nginx:alpine

# Копируем файл конфигурации Nginx
#COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

# Открываем порт для Nginx
EXPOSE 80

# Nginx запускается по умолчанию
