import { motion } from 'framer-motion';
import { AlertTriangle, Home, RefreshCw } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Button } from '../Button';
import { ROUTES } from '@/shared/constants';

interface ErrorFallbackProps {
  error: Error;
  onReset?: () => void;
}

/**
 * Компонент для отображения ошибки при падении React компонента
 * Используется как fallback UI в ErrorBoundary
 */
export function ErrorFallback({ error, onReset }: ErrorFallbackProps) {
  const isDevelopment = import.meta.env.DEV;

  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="max-w-md w-full text-center"
      >
        {/* Иконка ошибки */}
        <motion.div
          initial={{ scale: 0.8 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.1, type: 'spring', stiffness: 200 }}
          className="w-20 h-20 bg-error/20 rounded-full flex items-center justify-center mx-auto mb-6"
        >
          <AlertTriangle className="text-error" size={40} />
        </motion.div>

        {/* Заголовок */}
        <h1 className="text-3xl font-bold text-neutral-100 mb-2">
          Что-то пошло не так
        </h1>

        {/* Описание */}
        <p className="text-neutral-400 mb-6">
          Произошла непредвиденная ошибка. Мы уже работаем над её исправлением.
        </p>

        {/* Детали ошибки (только в режиме разработки) */}
        {isDevelopment && error && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="mb-6 p-4 bg-neutral-900 border border-neutral-800 rounded-lg text-left"
          >
            <p className="text-sm font-medium text-neutral-300 mb-2">
              Детали ошибки (только для разработки):
            </p>
            <pre className="text-xs text-neutral-400 overflow-auto max-h-40 font-mono">
              {error.toString()}
              {error.stack && `\n\n${error.stack}`}
            </pre>
          </motion.div>
        )}

        {/* Кнопки действий */}
        <div className="flex gap-3 justify-center flex-wrap">
          {onReset && (
            <Button
              variant="secondary"
              onClick={onReset}
              leftIcon={<RefreshCw size={18} />}
            >
              Попробовать снова
            </Button>
          )}
          <Link to={ROUTES.HOME}>
            <Button leftIcon={<Home size={18} />}>
              На главную
            </Button>
          </Link>
        </div>

        {/* Дополнительная информация */}
        <p className="text-xs text-neutral-500 mt-6">
          Если проблема повторяется, обратитесь в поддержку.
        </p>
      </motion.div>
    </div>
  );
}

