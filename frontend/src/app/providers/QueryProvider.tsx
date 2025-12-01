import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useState } from 'react';

interface QueryProviderProps {
  children: React.ReactNode;
}

/**
 * QueryProvider - настройка TanStack Query
 *
 * Конфигурация:
 * - staleTime: 30 секунд (данные считаются свежими)
 * - gcTime: 5 минут (время хранения в кеше)
 * - retry: 1 попытка при ошибке
 * - refetchOnWindowFocus: отключено
 */
export function QueryProvider({ children }: QueryProviderProps) {
  // Создаём QueryClient в useState чтобы избежать пересоздания при ререндерах
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30 * 1000, // 30 секунд
            gcTime: 5 * 60 * 1000, // 5 минут
            retry: 1,
            refetchOnWindowFocus: false,
          },
          mutations: {
            retry: 0,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} position="bottom" />
    </QueryClientProvider>
  );
}

