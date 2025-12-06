import { ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

interface AdminPageLayoutProps {
  title: string;
  description?: string;
  children: ReactNode;
  showBackButton?: boolean;
}

/**
 * Общий layout для страниц админ-панели
 */
export function AdminPageLayout({
  title,
  description,
  children,
  showBackButton = true,
}: AdminPageLayoutProps) {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex-1">
          {showBackButton && (
            <Link to={ROUTES.ADMIN}>
              <Button
                variant="ghost"
                size="sm"
                leftIcon={<ArrowLeft size={16} />}
                className="mb-4"
              >
                Назад
              </Button>
            </Link>
          )}
          <h1 className="text-2xl font-bold text-neutral-100">{title}</h1>
          {description && (
            <p className="text-sm text-neutral-400 mt-2">{description}</p>
          )}
        </div>
      </div>

      {/* Content */}
      {children}
    </div>
  );
}

