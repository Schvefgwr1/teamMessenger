import { Link } from 'react-router-dom';
import { Menu } from 'lucide-react';
import { cn } from '@/shared/lib/cn';
import { ROUTES } from '@/shared/constants';
import { useSidebarStore } from '@/widgets/Sidebar';
import { UserMenu } from './UserMenu';

/**
 * Header для основного приложения
 * Логотип слева, UserMenu справа
 */
export function Header() {
  const { toggleMobile } = useSidebarStore();

  return (
    <header
      className={cn(
        'fixed top-0 left-0 right-0 h-16 z-50',
        'bg-neutral-900/95 backdrop-blur-sm border-b border-neutral-800',
        'flex items-center justify-between px-4 md:px-6'
      )}
    >
      {/* Left side */}
      <div className="flex items-center gap-4">
        {/* Mobile menu button */}
        <button
          onClick={toggleMobile}
          className="lg:hidden p-2 rounded-lg text-neutral-400 hover:text-neutral-100 hover:bg-neutral-800 transition-colors"
          aria-label="Toggle menu"
        >
          <Menu size={20} />
        </button>

        {/* Logo */}
        <Link
          to={ROUTES.HOME}
          className="flex items-center gap-3 hover:opacity-80 transition-opacity"
        >
          <div className="w-9 h-9 bg-gradient-to-br from-primary-400 to-primary-600 rounded-lg flex items-center justify-center shadow-md shadow-primary-500/20">
            <span className="text-lg font-bold text-white">TM</span>
          </div>
          <span className="text-lg font-semibold text-neutral-100 hidden sm:block">
            Team Messenger
          </span>
        </Link>
      </div>

      {/* Right side */}
      <div className="flex items-center gap-3">
        {/* TODO: Notifications */}
        {/* TODO: Search */}

        {/* User menu */}
        <UserMenu />
      </div>
    </header>
  );
}

