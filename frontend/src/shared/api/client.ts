import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8090';

/**
 * Формируем baseURL для API
 * Если VITE_API_URL начинается с /, это относительный путь (production через nginx)
 * Если начинается с http, это абсолютный URL (development)
 */
const getApiBaseURL = () => {
  if (API_BASE_URL.startsWith('/')) {
    // Относительный путь (production) - просто добавляем /api/v1
    return `${API_BASE_URL}/v1`;
  } else {
    // Абсолютный URL (development) - добавляем /api/v1
    return `${API_BASE_URL}/api/v1`;
  }
};

/**
 * Axios instance для работы с API
 */
export const apiClient = axios.create({
  baseURL: getApiBaseURL(),
  timeout: 30000,
  withCredentials: true, // Для CORS с credentials
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Получить токен из localStorage
 * Используем функцию вместо прямого импорта store для избежания circular dependency
 */
function getToken(): string | null {
  try {
    const authStorage = localStorage.getItem('auth-storage');
    if (authStorage) {
      const parsed = JSON.parse(authStorage);
      return parsed.state?.token || null;
    }
  } catch {
    // ignore parse errors
  }
  return null;
}

/**
 * Очистить токен и перенаправить на логин
 */
function handleUnauthorized(): void {
  localStorage.removeItem('auth-storage');
  // Редирект только если не на странице логина
  if (!window.location.pathname.includes('/login')) {
    window.location.href = '/login';
  }
}

// Request interceptor - добавление токена
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - обработка ошибок
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      handleUnauthorized();
    }
    if (error.response?.status === 429) {
      console.error('Too many requests - rate limit exceeded');
    }
    return Promise.reject(error);
  }
);

export { API_BASE_URL };

