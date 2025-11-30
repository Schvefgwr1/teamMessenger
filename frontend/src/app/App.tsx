import { ToastProvider, Button, Card, Avatar, Badge, Input, Spinner, toast } from '@/shared/ui';

function App() {
  const handleToast = () => {
    toast.success('Успешно!', { description: 'Это демонстрация toast уведомления' });
  };

  return (
    <>
      <ToastProvider />
      <div className="min-h-screen bg-neutral-950 text-neutral-100 p-8">
        <div className="max-w-4xl mx-auto space-y-8">
          <header className="text-center">
            <h1 className="text-4xl font-bold mb-2">Team Messenger</h1>
            <p className="text-neutral-400">UI Kit Demo</p>
          </header>

          {/* Buttons */}
          <Card>
            <Card.Header>
              <Card.Title>Buttons</Card.Title>
            </Card.Header>
            <Card.Body className="flex flex-wrap gap-3">
              <Button variant="primary">Primary</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="outline">Outline</Button>
              <Button variant="ghost">Ghost</Button>
              <Button variant="danger">Danger</Button>
              <Button variant="primary" isLoading>
                Loading
              </Button>
            </Card.Body>
          </Card>

          {/* Inputs */}
          <Card>
            <Card.Header>
              <Card.Title>Inputs</Card.Title>
            </Card.Header>
            <Card.Body className="space-y-4">
              <Input label="Email" placeholder="user@example.com" type="email" />
              <Input
                label="Пароль"
                placeholder="••••••••"
                type="password"
                hint="Минимум 8 символов"
              />
              <Input label="С ошибкой" placeholder="Введите данные" error="Это поле обязательно" />
            </Card.Body>
          </Card>

          {/* Avatars & Badges */}
          <Card>
            <Card.Header>
              <Card.Title>Avatars & Badges</Card.Title>
            </Card.Header>
            <Card.Body className="flex flex-wrap items-center gap-4">
              <Avatar fallback="JD" size="xs" />
              <Avatar fallback="AB" size="sm" status="online" />
              <Avatar fallback="XY" size="md" status="away" />
              <Avatar fallback="МИ" size="lg" status="busy" />
              <Avatar fallback="ПР" size="xl" status="offline" />

              <div className="flex gap-2 ml-4">
                <Badge>Default</Badge>
                <Badge variant="primary">Primary</Badge>
                <Badge variant="success">Success</Badge>
                <Badge variant="warning">Warning</Badge>
                <Badge variant="error">Error</Badge>
              </div>
            </Card.Body>
          </Card>

          {/* Toast Demo */}
          <Card>
            <Card.Header>
              <Card.Title>Toast Notifications</Card.Title>
            </Card.Header>
            <Card.Body className="flex flex-wrap gap-3">
              <Button onClick={handleToast}>Show Toast</Button>
              <Button variant="secondary" onClick={() => toast.error('Ошибка!')}>
                Error Toast
              </Button>
              <Button variant="secondary" onClick={() => toast.warning('Внимание!')}>
                Warning Toast
              </Button>
              <Button variant="secondary" onClick={() => toast.info('Информация')}>
                Info Toast
              </Button>
            </Card.Body>
          </Card>

          {/* Spinner */}
          <Card>
            <Card.Header>
              <Card.Title>Spinners</Card.Title>
            </Card.Header>
            <Card.Body className="flex items-center gap-6">
              <Spinner size="sm" />
              <Spinner size="md" />
              <Spinner size="lg" />
            </Card.Body>
          </Card>
        </div>
      </div>
    </>
  );
}

export default App;
