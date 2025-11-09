# ---------- Build stage ----------
FROM node:20-alpine AS build
WORKDIR /app

# Устанавливаем зависимости
COPY package*.json ./
RUN npm ci

# Копируем исходники
COPY . .

# Сборка production-бандла
RUN npm run build

# ---------- Runtime stage ----------
FROM nginx:1.27-alpine

# Конфиг nginx с fallback на index.html (SPA)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Статика фронта
COPY --from=build /app/dist/ /usr/share/nginx/html/

# Nginx слушает 80 внутри контейнера
EXPOSE 80

# (опционально) healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
