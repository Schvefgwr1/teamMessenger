import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { TaskBoard } from '@/widgets/TaskBoard';
import { CreateTaskModal } from '@/features/task/create-task';

/**
 * Страница со списком задач (Kanban доска)
 */
export function TasksPage() {
  const navigate = useNavigate();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const handleTaskClick = (taskId: number) => {
    navigate(`/tasks/${taskId}`);
  };

  return (
    <div className="flex flex-col min-h-0">
      <TaskBoard
        onTaskClick={handleTaskClick}
        onCreateTask={() => setIsCreateModalOpen(true)}
      />

      <CreateTaskModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
      />
    </div>
  );
}

