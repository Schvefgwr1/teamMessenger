import { Component, ErrorInfo, ReactNode } from 'react';
import { ErrorFallback } from './ErrorFallback';

export interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

/**
 * Error Boundary для перехвата ошибок React компонентов
 * 
 * @example
 * ```tsx
 * <ErrorBoundary>
 *   <App />
 * </ErrorBoundary>
 * ```
 * 
 * @example
 * ```tsx
 * <ErrorBoundary
 *   fallback={<CustomErrorUI />}
 *   onError={(error, errorInfo) => {
 *     // Отправка в error tracking service
 *     logErrorToService(error, errorInfo);
 *   }}
 * >
 *   <App />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Логируем ошибку в консоль для разработки
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // Вызываем пользовательский обработчик если предоставлен
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // TODO: Здесь можно добавить отправку в error tracking service
    // Например: Sentry, LogRocket, etc.
    // logErrorToService(error, errorInfo);
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
    });
  };

  render() {
    if (this.state.hasError && this.state.error) {
      // Если предоставлен кастомный fallback, используем его
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Иначе используем дефолтный ErrorFallback
      return (
        <ErrorFallback error={this.state.error} onReset={this.handleReset} />
      );
    }

    return this.props.children;
  }
}

