import { motion, AnimatePresence } from 'framer-motion';
import { useLocation } from 'react-router-dom';
import { ReactNode } from 'react';
import { pageTransition } from '@/shared/lib/animations';

interface PageTransitionProps {
  children: ReactNode;
}

/**
 * Компонент для анимации переходов между страницами
 * Обертывает содержимое страницы для плавных переходов
 * 
 * @example
 * ```tsx
 * function DashboardPage() {
 *   return (
 *     <PageTransition>
 *       <div>Content</div>
 *     </PageTransition>
 *   );
 * }
 * ```
 */
export function PageTransition({ children }: PageTransitionProps) {
  const location = useLocation();

  return (
    <AnimatePresence mode="wait">
      <motion.div
        key={location.pathname}
        variants={pageTransition}
        initial="initial"
        animate="animate"
        exit="exit"
        className="w-full h-full"
      >
        {children}
      </motion.div>
    </AnimatePresence>
  );
}

