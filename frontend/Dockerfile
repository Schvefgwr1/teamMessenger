# ===============================
# Stage 1: Build
# ===============================
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем package files
COPY package.json package-lock.json ./

# Устанавливаем зависимости
RUN npm ci

# Копируем исходный код
COPY . .

# Build arguments для env переменных
# Frontend обращается к API через backend nginx на localhost:8090
# В браузере пользователя запросы идут напрямую на backend nginx
ARG VITE_API_URL=http://localhost:8090
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Собираем приложение
RUN npm run build

# ===============================
# Stage 2: Production
# ===============================
FROM nginx:alpine AS production

# Устанавливаем wget для healthcheck
RUN apk add --no-cache wget

# Копируем собранное приложение
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose порт
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

