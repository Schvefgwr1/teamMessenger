import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryProvider, AuthProvider } from './providers';
import { ProtectedRoute, AdminRoute, GuestRoute } from './router';
import { ToastProvider, ErrorBoundary } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

// Layouts
import { AuthLayout, MainLayout, AdminLayout } from '@/widgets/layouts';

// Pages
import { LoginPage } from '@/pages/auth/LoginPage';
import { RegisterPage } from '@/pages/auth/RegisterPage';
import { DashboardPage } from '@/pages/dashboard/DashboardPage';
import { ProfilePage } from '@/pages/profile/ProfilePage';
import { ChatsPage, ChatDetailPage } from '@/pages/chats';
import { TasksPage, TaskDetailPage } from '@/pages/tasks';
import { NotFoundPage, ForbiddenPage } from '@/pages/errors';
import {
  AdminDashboard,
  RolesPage,
  PermissionsPage,
  ChangeUserRolePage,
  ChatRolesPage,
  ChatPermissionsPage,
  TaskStatusesPage,
} from '@/pages/admin';

function App() {
  return (
    <ErrorBoundary>
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
                  <Route path={ROUTES.PROFILE} element={<ProfilePage />} />
                  <Route path={ROUTES.CHATS} element={<ChatsPage />}>
                    <Route path=":chatId" element={<ChatDetailPage />} />
                  </Route>
                  <Route path={ROUTES.TASKS} element={<TasksPage />} />
                  <Route path="/tasks/:taskId" element={<TaskDetailPage />} />
                </Route>
              </Route>

              {/* Admin routes (требуют admin права) */}
              <Route element={<AdminRoute />}>
                <Route element={<AdminLayout />}>
                  <Route path={ROUTES.ADMIN} element={<AdminDashboard />} />
                  <Route path={ROUTES.ADMIN_ROLES} element={<RolesPage />} />
                  <Route path={ROUTES.ADMIN_PERMISSIONS} element={<PermissionsPage />} />
                  <Route path={ROUTES.ADMIN_CHANGE_USER_ROLE} element={<ChangeUserRolePage />} />
                  <Route path={ROUTES.ADMIN_CHAT_ROLES} element={<ChatRolesPage />} />
                  <Route path={ROUTES.ADMIN_CHAT_PERMISSIONS} element={<ChatPermissionsPage />} />
                  <Route path={ROUTES.ADMIN_TASK_STATUSES} element={<TaskStatusesPage />} />
                </Route>
              </Route>

              {/* Error pages */}
              <Route path={ROUTES.FORBIDDEN} element={<ForbiddenPage />} />
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </QueryProvider>
    </ErrorBoundary>
  );
}

export default App;
