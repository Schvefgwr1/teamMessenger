import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import { cn } from '@/shared/lib';

interface DropdownProps {
  children: React.ReactNode;
}

export function Dropdown({ children }: DropdownProps) {
  return <DropdownMenu.Root>{children}</DropdownMenu.Root>;
}

Dropdown.Trigger = DropdownMenu.Trigger;

Dropdown.Content = function DropdownContent({
  children,
  className,
  align = 'end',
  sideOffset = 8,
}: {
  children: React.ReactNode;
  className?: string;
  align?: 'start' | 'center' | 'end';
  sideOffset?: number;
}) {
  return (
    <DropdownMenu.Portal>
      <DropdownMenu.Content
        align={align}
        sideOffset={sideOffset}
        className={cn(
          'z-50 min-w-[180px] overflow-hidden rounded-xl p-1',
          'bg-neutral-900 border border-neutral-800 shadow-lg',
          'data-[state=open]:animate-in data-[state=closed]:animate-out',
          'data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
          'data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95',
          'data-[side=bottom]:slide-in-from-top-2',
          'data-[side=left]:slide-in-from-right-2',
          'data-[side=right]:slide-in-from-left-2',
          'data-[side=top]:slide-in-from-bottom-2',
          className
        )}
      >
        {children}
      </DropdownMenu.Content>
    </DropdownMenu.Portal>
  );
};

Dropdown.Item = function DropdownItem({
  children,
  className,
  destructive = false,
  ...props
}: DropdownMenu.DropdownMenuItemProps & { destructive?: boolean }) {
  return (
    <DropdownMenu.Item
      className={cn(
        'flex items-center gap-2 px-3 py-2 text-sm rounded-lg cursor-pointer',
        'outline-none transition-colors',
        destructive
          ? 'text-error focus:bg-error/20'
          : 'text-neutral-300 focus:bg-neutral-800 focus:text-neutral-100',
        className
      )}
      {...props}
    >
      {children}
    </DropdownMenu.Item>
  );
};

Dropdown.Separator = function DropdownSeparator({ className }: { className?: string }) {
  return <DropdownMenu.Separator className={cn('h-px my-1 bg-neutral-800', className)} />;
};

Dropdown.Label = function DropdownLabel({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <DropdownMenu.Label className={cn('px-3 py-1.5 text-xs text-neutral-500', className)}>
      {children}
    </DropdownMenu.Label>
  );
};

