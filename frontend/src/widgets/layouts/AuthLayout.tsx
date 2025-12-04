import {Outlet} from 'react-router-dom';
import {motion} from 'framer-motion';

/**
 * Layout для страниц авторизации (login, register)
 * Брендинг слева, форма справа
 */
export function AuthLayout() {
    return (
        <div className="min-h-screen bg-neutral-950 flex">
            {/* Left side - Branding (скрыт на мобильных) */}
            <div
                className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-neutral-900 via-neutral-950 to-neutral-900 items-center justify-center p-12 relative overflow-hidden">
                {/* Background decoration */}
                <div className="absolute inset-0 opacity-30">
                    <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary-500/20 rounded-full blur-3xl"/>
                    <div className="absolute bottom-1/4 right-1/4 w-64 h-64 bg-primary-600/10 rounded-full blur-3xl"/>
                </div>

                <motion.div
                    initial={{opacity: 0, y: 20}}
                    animate={{opacity: 1, y: 0}}
                    transition={{duration: 0.6}}
                    className="max-w-md text-center relative z-10"
                >
                    {/* Logo */}
                    <div
                        className="w-20 h-20 bg-gradient-to-br from-primary-400 to-primary-600 rounded-2xl mx-auto mb-8 flex items-center justify-center shadow-lg shadow-primary-500/25">
                        <span className="text-3xl font-bold text-white">TM</span>
                    </div>

                    <h1 className="text-4xl font-bold text-neutral-100 mb-4">
                        Team Messenger
                    </h1>
                    <p className="text-neutral-400 text-lg leading-relaxed">
                        Общайтесь с командой, управляйте задачами и работайте вместе
                        эффективнее
                    </p>

                    {/* Features */}
                    <div className="mt-12 space-y-4 text-left">
                        <Feature
                            icon="💬"
                            title="Групповые чаты"
                            description="Создавайте чаты для команд и проектов"
                        />
                        <Feature
                            icon="✅"
                            title="Управление задачами"
                            description="Kanban-доска для отслеживания прогресса"
                        />
                        <Feature
                            icon="🔒"
                            title="Гибкие права"
                            description="Настраиваемые роли и разрешения"
                        />
                    </div>
                </motion.div>
            </div>

            {/* Right side - Auth forms */}
            <div className="flex-1 flex justify-center p-6 overflow-y-auto max-h-screen">
                <Outlet/>
            </div>
        </div>
    );
}

interface FeatureProps {
    icon: string;
    title: string;
    description: string;
}

function Feature({icon, title, description}: FeatureProps) {
    return (
        <div className="flex items-start gap-4 p-4 rounded-xl bg-neutral-800/30 backdrop-blur-sm">
            <div className="flex" style={{minHeight: "-webkit-fill-available", alignItems: "center"}}>
                <span className="text-2xl">{icon}</span>
            </div>
            <div>
                <h3 className="font-medium text-neutral-200">{title}</h3>
                <p className="text-sm text-neutral-500">{description}</p>
            </div>
        </div>
    );
}

