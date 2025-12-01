import { Link } from 'react-router-dom';
import { Button } from '@/shared/ui';
import { Home, ArrowLeft, ShieldAlert } from 'lucide-react';
import { ROUTES } from '@/shared/constants';

export function ForbiddenPage() {
  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4">
      <div className="text-center">
        <div className="w-20 h-20 bg-error/20 rounded-full flex items-center justify-center mx-auto mb-6">
          <ShieldAlert className="text-error" size={40} />
        </div>
        <h1 className="text-6xl font-bold text-neutral-800">403</h1>
        <h2 className="text-2xl font-semibold text-neutral-100 mt-4">
          Доступ запрещён
        </h2>
        <p className="text-neutral-400 mt-2 max-w-md">
          У вас нет прав для просмотра этой страницы.
          Обратитесь к администратору, если считаете, что это ошибка.
        </p>
        <div className="flex gap-4 justify-center mt-8">
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
        </div>
      </div>
    </div>
  );
}

