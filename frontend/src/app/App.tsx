import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryProvider, AuthProvider } from './providers';
import { ProtectedRoute, AdminRoute, GuestRoute } from './router';
import { ToastProvider } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

// Layouts
import { AuthLayout, MainLayout, AdminLayout } from '@/widgets/layouts';

// Pages
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
              <Route element={<AuthLayout />}>
                <Route path={ROUTES.LOGIN} element={<LoginPage />} />
                <Route path={ROUTES.REGISTER} element={<RegisterPage />} />
              </Route>
            </Route>

            {/* Protected routes (требуют авторизации) */}
            <Route element={<ProtectedRoute />}>
              <Route element={<MainLayout />}>
                <Route path={ROUTES.HOME} element={<DashboardPage />} />
                <Route path={ROUTES.PROFILE} element={<div className="text-neutral-100">ProfilePage (TODO)</div>} />
                <Route path={ROUTES.CHATS} element={<div className="text-neutral-100">ChatsPage (TODO)</div>} />
                <Route path="/chats/:chatId" element={<div className="text-neutral-100">ChatDetailPage (TODO)</div>} />
                <Route path={ROUTES.TASKS} element={<div className="text-neutral-100">TasksPage (TODO)</div>} />
                <Route path="/tasks/:taskId" element={<div className="text-neutral-100">TaskDetailPage (TODO)</div>} />
              </Route>
            </Route>

            {/* Admin routes (требуют admin права) */}
            <Route element={<AdminRoute />}>
              <Route element={<AdminLayout />}>
                <Route path={ROUTES.ADMIN} element={<div className="text-neutral-100">AdminDashboard (TODO)</div>} />
                <Route path={ROUTES.ADMIN_ROLES} element={<div className="text-neutral-100">RolesPage (TODO)</div>} />
                <Route path={ROUTES.ADMIN_PERMISSIONS} element={<div className="text-neutral-100">PermissionsPage (TODO)</div>} />
                <Route path={ROUTES.ADMIN_CHAT_ROLES} element={<div className="text-neutral-100">ChatRolesPage (TODO)</div>} />
                <Route path={ROUTES.ADMIN_CHAT_PERMISSIONS} element={<div className="text-neutral-100">ChatPermissionsPage (TODO)</div>} />
                <Route path={ROUTES.ADMIN_TASK_STATUSES} element={<div className="text-neutral-100">TaskStatusesPage (TODO)</div>} />
              </Route>
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
