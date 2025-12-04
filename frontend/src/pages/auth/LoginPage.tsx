import {useState} from 'react';
import {Link, useLocation, useNavigate} from 'react-router-dom';
import {useMutation} from '@tanstack/react-query';
import {authApi, useAuthStore} from '@/entities/session';
import {transformUserResponse, userApi} from '@/entities/user';
import {Button, Card, Input, toast} from '@/shared/ui';
import {ROUTES} from '@/shared/constants';
import {Eye, EyeOff, Lock, Mail} from 'lucide-react';
import {motion} from "framer-motion";

export function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const {login, setToken} = useAuthStore();

    const [formData, setFormData] = useState({
        login: '',
        password: '',
    });
    const [showPassword, setShowPassword] = useState(false);

    // URL для редиректа после логина
    const from = location.state?.from?.pathname || ROUTES.HOME;

    const loginMutation = useMutation({
        mutationFn: async (data: { login: string; password: string }) => {
            // 1. Получаем токен
            const authResponse = await authApi.login(data);
            const {token} = authResponse.data;

            // Временно сохраняем токен для следующего запроса
            setToken(token);

            // 2. Загружаем данные пользователя
            const userResponse = await userApi.getMe();
            const user = transformUserResponse(userResponse.data);

            return {token, user};
        },
        onSuccess: ({token, user}) => {
            login(token, user);
            toast.success('Добро пожаловать!');
            navigate(from, {replace: true});
        },
        onError: (error: Error & { response?: { data?: { error?: string } } }) => {
            const message = error.response?.data?.error || 'Ошибка входа';
            toast.error(message);
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.login || !formData.password) {
            toast.error('Заполните все поля');
            return;
        }
        loginMutation.mutate(formData);
    };

    return (
        <motion.div
            initial={{opacity: 0, x: 20}}
            animate={{opacity: 1, x: 0}}
            transition={{duration: 0.4, delay: 0.2}}
            className="w-full max-w-md flex justify-center flex-col"
        >

            {/* Logo (для мобильных, на desktop показывается в AuthLayout) */}
            <div className="text-center mb-8 lg:hidden">
                <div
                    className="w-16 h-16 bg-gradient-to-br from-primary-400 to-primary-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
                    <span className="text-2xl font-bold text-white">TM</span>
                </div>
            </div>

            <div className="text-center mb-6">
                <h1 className="text-2xl font-bold text-neutral-100">Вход в систему</h1>
                <p className="text-neutral-400 mt-2">Введите данные для входа</p>
            </div>

            <Card>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <Input
                        label="Логин"
                        placeholder="username"
                        value={formData.login}
                        onChange={(e) => setFormData({...formData, login: e.target.value})}
                        leftIcon={<Mail size={18}/>}
                        autoComplete="username"
                    />

                    <Input
                        label="Пароль"
                        type={showPassword ? 'text' : 'password'}
                        placeholder="••••••••"
                        value={formData.password}
                        onChange={(e) => setFormData({...formData, password: e.target.value})}
                        leftIcon={<Lock size={18}/>}
                        rightIcon={
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="hover:text-neutral-300 transition-colors"
                            >
                                {showPassword ? <EyeOff size={18}/> : <Eye size={18}/>}
                            </button>
                        }
                        autoComplete="current-password"
                    />

                    <Button
                        type="submit"
                        className="w-full"
                        isLoading={loginMutation.isPending}
                    >
                        Войти
                    </Button>
                </form>

                <div className="mt-6 text-center text-sm">
                    <span className="text-neutral-400">Нет аккаунта? </span>
                    <Link
                        to={ROUTES.REGISTER}
                        className="text-primary-400 hover:text-primary-300 font-medium"
                    >
                        Зарегистрироваться
                    </Link>
                </div>
            </Card>
        </motion.div>
    );
}
