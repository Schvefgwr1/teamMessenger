import * as Dialog from '@radix-ui/react-dialog';
import { X } from 'lucide-react';
import { cn } from '@/shared/lib';
import { Button } from '../Button';

interface ModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
}

export function Modal({ open, onOpenChange, children }: ModalProps) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      {children}
    </Dialog.Root>
  );
}

Modal.Trigger = Dialog.Trigger;

Modal.Content = function ModalContent({
  children,
  className,
  title,
  description,
  showClose = true,
}: {
  children: React.ReactNode;
  className?: string;
  title?: string;
  description?: string;
  showClose?: boolean;
}) {
  return (
    <Dialog.Portal>
      <Dialog.Overlay className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0" />
      <Dialog.Content
        className={cn(
          'fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50',
          'w-full max-w-lg max-h-[85vh] overflow-auto',
          'bg-neutral-900 border border-neutral-800 rounded-xl shadow-xl',
          'p-6',
          'data-[state=open]:animate-in data-[state=closed]:animate-out',
          'data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
          'data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95',
          'data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%]',
          'data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%]',
          'duration-200',
          className
        )}
      >
        {(title || showClose) && (
          <div className="flex items-center justify-between mb-4">
            <div>
              {title && (
                <Dialog.Title className="text-lg font-semibold text-neutral-100">
                  {title}
                </Dialog.Title>
              )}
              {description && (
                <Dialog.Description className="text-sm text-neutral-400 mt-1">
                  {description}
                </Dialog.Description>
              )}
            </div>
            {showClose && (
              <Dialog.Close asChild>
                <Button variant="ghost" size="icon-sm" className="text-neutral-500">
                  <X className="w-4 h-4" />
                </Button>
              </Dialog.Close>
            )}
          </div>
        )}
        {children}
      </Dialog.Content>
    </Dialog.Portal>
  );
};

Modal.Footer = function ModalFooter({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div className={cn('flex justify-end gap-2 mt-6 pt-4 border-t border-neutral-800', className)}>
      {children}
    </div>
  );
};

Modal.Close = Dialog.Close;

