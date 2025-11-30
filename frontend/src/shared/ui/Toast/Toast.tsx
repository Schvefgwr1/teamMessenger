import { Toaster as SonnerToaster, toast as sonnerToast } from 'sonner';

/**
 * Компонент Toast контейнер
 * Размещается в корне приложения
 */
export function ToastProvider() {
  return (
    <SonnerToaster
      position="top-right"
      toastOptions={{
        classNames: {
          toast:
            'bg-neutral-900 border border-neutral-800 text-neutral-100 shadow-lg rounded-xl',
          description: 'text-neutral-400',
          actionButton: 'bg-primary-500 text-white',
          cancelButton: 'bg-neutral-800 text-neutral-100',
          error: 'bg-error/10 border-error text-error',
          success: 'bg-success/10 border-success text-success',
          warning: 'bg-warning/10 border-warning text-warning',
          info: 'bg-info/10 border-info text-info',
        },
      }}
      closeButton
      richColors
      expand
    />
  );
}

/**
 * Утилита для показа toast уведомлений
 */
export const toast = {
  success: (message: string, options?: { description?: string }) =>
    sonnerToast.success(message, options),

  error: (message: string, options?: { description?: string }) =>
    sonnerToast.error(message, options),

  warning: (message: string, options?: { description?: string }) =>
    sonnerToast.warning(message, options),

  info: (message: string, options?: { description?: string }) =>
    sonnerToast.info(message, options),

  loading: (message: string) => sonnerToast.loading(message),

  dismiss: (toastId?: string | number) => sonnerToast.dismiss(toastId),

  promise: <T,>(
    promise: Promise<T>,
    options: {
      loading: string;
      success: string | ((data: T) => string);
      error: string | ((error: Error) => string);
    }
  ) => sonnerToast.promise(promise, options),
};

