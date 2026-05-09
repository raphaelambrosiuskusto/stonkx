# ЭТАП 1: Сборка (Node.js)
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем конфиги и зависимости
COPY package*.json ./
RUN npm install

# Копируем всё остальное
COPY . .

# Собираем JS (используем прямой вызов esbuild, раз в scripts пусто)
RUN ./node_modules/.bin/esbuild src/app.ts --bundle --outfile=dist/app.js

# ВАЖНО: убеждаемся, что в папке dist есть всё для работы сайта
# Если index.html и css/ лежат в корне /app, копируем их в /app/dist
RUN cp index.html dist/ && cp -r css dist/ || true


# ЭТАП 2: Раздача (Nginx)
FROM nginx:stable-alpine

# Наш конфиг для проксирования запросов к Go
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем содержимое собранной папки dist в стандартную папку nginx
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
