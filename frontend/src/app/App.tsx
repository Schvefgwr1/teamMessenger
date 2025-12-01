import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryProvider, AuthProvider } from './providers';
import { ProtectedRoute, AdminRoute, GuestRoute } from './router';
import { ToastProvider } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

// Pages (временные заглушки)
import { LoginPage } from '@/pages/auth/LoginPage';
import { RegisterPage } from '@/pages/auth/RegisterPage';
import { DashboardPage } from '@/pages/dashboard/DashboardPage';
import { NotFoundPage } from '@/pages/errors/NotFoundPage';
import { ForbiddenPage } from '@/pages/errors/ForbiddenPage';

function App() {
  return (
    <QueryProvider>
      <AuthProvider>
        <BrowserRouter>
          <ToastProvider />
          <Routes>
            {/* Guest routes (только для неавторизованных) */}
            <Route element={<GuestRoute />}>
              <Route path={ROUTES.LOGIN} element={<LoginPage />} />
              <Route path={ROUTES.REGISTER} element={<RegisterPage />} />
            </Route>

            {/* Protected routes (требуют авторизации) */}
            <Route element={<ProtectedRoute />}>
              <Route path={ROUTES.HOME} element={<DashboardPage />} />
              {/* TODO: добавить остальные защищённые маршруты */}
            </Route>

            {/* Admin routes (требуют admin права) */}
            <Route element={<AdminRoute />}>
              <Route path={ROUTES.ADMIN} element={<div>Admin Dashboard</div>} />
              {/* TODO: добавить админ маршруты */}
            </Route>

            {/* Error pages */}
            <Route path={ROUTES.FORBIDDEN} element={<ForbiddenPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryProvider>
  );
}

export default App;
