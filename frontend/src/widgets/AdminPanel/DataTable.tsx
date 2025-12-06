import { ReactNode } from 'react';
import { Card, Skeleton } from '@/shared/ui';
import { cn } from '@/shared/lib';

interface Column<T> {
  key: string;
  header: string;
  render: (item: T) => ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  data: T[];
  columns: Column<T>[];
  isLoading?: boolean;
  emptyMessage?: string;
  onRowClick?: (item: T) => void;
  className?: string;
}

/**
 * Универсальная таблица данных для админ-панели
 */
export function DataTable<T extends { id?: number | string }>({
  data,
  columns,
  isLoading = false,
  emptyMessage = 'Нет данных',
  onRowClick,
  className,
}: DataTableProps<T>) {
  if (isLoading) {
    return (
      <Card>
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="flex gap-4">
              {columns.map((_, colIdx) => (
                <Skeleton key={colIdx} className="h-10 flex-1" />
              ))}
            </div>
          ))}
        </div>
      </Card>
    );
  }

  if (data.length === 0) {
    return (
      <Card>
        <div className="text-center py-12">
          <p className="text-neutral-500">{emptyMessage}</p>
        </div>
      </Card>
    );
  }

  return (
    <Card className={cn('overflow-hidden', className)}>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-neutral-800">
              {columns.map((column) => (
                <th
                  key={column.key}
                  className={cn(
                    'px-4 py-3 text-left text-sm font-semibold text-neutral-300',
                    column.className
                  )}
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((item, index) => (
              <tr
                key={item.id ?? index}
                onClick={() => onRowClick?.(item)}
                className={cn(
                  'border-b border-neutral-800/50 transition-colors',
                  onRowClick && 'cursor-pointer hover:bg-neutral-800/30'
                )}
              >
                {columns.map((column) => (
                  <td
                    key={column.key}
                    className={cn('px-4 py-3 text-sm text-neutral-200', column.className)}
                  >
                    {column.render(item)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
}

