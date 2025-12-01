import { Link } from 'react-router-dom';
import { Button } from '@/shared/ui';
import { Home, ArrowLeft } from 'lucide-react';
import { ROUTES } from '@/shared/constants';

export function NotFoundPage() {
  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4">
      <div className="text-center">
        <h1 className="text-9xl font-bold text-neutral-800">404</h1>
        <h2 className="text-2xl font-semibold text-neutral-100 mt-4">
          Страница не найдена
        </h2>
        <p className="text-neutral-400 mt-2 max-w-md">
          Запрашиваемая страница не существует или была удалена.
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

