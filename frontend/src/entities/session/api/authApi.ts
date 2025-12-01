import { apiClient } from '@/shared/api';

/**
 * Запрос на вход
 */
export interface LoginRequest {
  login: string;
  password: string;
}

/**
 * Ответ при успешном входе
 */
export interface LoginResponse {
  token: string;
  userID: string;
}

/**
 * Запрос на регистрацию (JSON часть в FormData.data)
 */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  description?: string;
  gender?: string;
  age?: number;
  roleID?: number; // Для dev режима, в production = 1
}

/**
 * ID роли по умолчанию (обычный пользователь)
 */
export const DEFAULT_ROLE_ID = 1;

/**
 * Флаг для показа поля roleID в форме регистрации (только для dev)
 */
export const SHOW_ROLE_FIELD_IN_REGISTER = false;

/**
 * Auth API endpoints
 */
export const authApi = {
  /**
   * Вход в систему
   * POST /api/v1/auth/login
   */
  login: (data: LoginRequest) =>
    apiClient.post<LoginResponse>('/auth/login', data),

  /**
   * Регистрация нового пользователя
   * POST /api/v1/auth/register (multipart/form-data)
   *
   * FormData:
   * - data: JSON string с RegisterRequest
   * - file: File (аватар, опционально)
   */
  register: (formData: FormData) =>
    apiClient.post('/auth/register', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Выход из системы
   * POST /api/v1/auth/logout
   */
  logout: () => apiClient.post('/auth/logout'),
};

/**
 * Создать FormData для регистрации
 */
export function createRegisterFormData(
  data: RegisterRequest,
  avatar?: File
): FormData {
  const formData = new FormData();
  
  // Подготавливаем данные для JSON, исключая пустые опциональные поля
  const jsonData: RegisterRequest = {
    username: data.username,
    email: data.email,
    password: data.password,
    roleID: data.roleID ?? DEFAULT_ROLE_ID,
  };
  
  // Добавляем опциональные поля только если они заполнены
  if (data.description && data.description.trim()) {
    jsonData.description = data.description.trim();
  }
  if (data.gender && data.gender.trim()) {
    jsonData.gender = data.gender.trim();
  }
  if (data.age !== undefined && data.age !== null) {
    jsonData.age = data.age;
  }
  
  formData.append('data', JSON.stringify(jsonData));
  
  if (avatar) {
    formData.append('file', avatar);
  }
  
  return formData;
}

