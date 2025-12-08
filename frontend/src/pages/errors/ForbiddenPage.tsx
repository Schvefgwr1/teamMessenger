import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import { Button } from '@/shared/ui';
import { Home, ArrowLeft, ShieldAlert } from 'lucide-react';
import { ROUTES } from '@/shared/constants';

export function ForbiddenPage() {
  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="text-center max-w-lg"
      >
        {/* Иконка */}
        <motion.div
          initial={{ scale: 0.8, rotate: -10 }}
          animate={{ scale: 1, rotate: 0 }}
          transition={{ delay: 0.1, type: 'spring', stiffness: 200 }}
          className="w-24 h-24 bg-error/20 rounded-full flex items-center justify-center mx-auto mb-6"
        >
          <ShieldAlert className="text-error" size={48} />
        </motion.div>

        {/* Код ошибки */}
        <motion.h1
          initial={{ opacity: 0, scale: 0.5 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
          className="text-8xl font-bold text-neutral-800 leading-none"
        >
          403
        </motion.h1>

        {/* Заголовок */}
        <motion.h2
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="text-2xl font-semibold text-neutral-100 mt-4 mb-2"
        >
          Доступ запрещён
        </motion.h2>

        {/* Описание */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4 }}
          className="text-neutral-400 mb-8 max-w-md mx-auto"
        >
          У вас нет прав для просмотра этой страницы.
          Обратитесь к администратору, если считаете, что это ошибка.
        </motion.p>

        {/* Кнопки действий */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="flex gap-4 justify-center flex-wrap"
        >
          <Button
            variant="secondary"
            onClick={() => window.history.back()}
            leftIcon={<ArrowLeft size={18} />}
          >
            Назад
          </Button>
          <Link to={ROUTES.HOME}>
            <Button leftIcon={<Home size={18} />}>
              На главную
            </Button>
          </Link>
        </motion.div>
      </motion.div>
    </div>
  );
}

