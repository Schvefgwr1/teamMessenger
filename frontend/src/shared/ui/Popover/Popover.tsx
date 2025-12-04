import * as PopoverPrimitive from '@radix-ui/react-popover';
import { cn } from '@/shared/lib';

interface PopoverProps {
  children: React.ReactNode;
}

export function Popover({ children }: PopoverProps) {
  return <PopoverPrimitive.Root>{children}</PopoverPrimitive.Root>;
}

Popover.Trigger = PopoverPrimitive.Trigger;

Popover.Content = function PopoverContent({
  children,
  className,
  align = 'center',
  side = 'top',
  sideOffset = 8,
}: {
  children: React.ReactNode;
  className?: string;
  align?: 'start' | 'center' | 'end';
  side?: 'top' | 'right' | 'bottom' | 'left';
  sideOffset?: number;
}) {
  return (
    <PopoverPrimitive.Portal>
      <PopoverPrimitive.Content
        align={align}
        side={side}
        sideOffset={sideOffset}
        className={cn(
          'z-50 w-64 rounded-xl p-4',
          'bg-neutral-900 border border-neutral-800 shadow-xl',
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
      </PopoverPrimitive.Content>
    </PopoverPrimitive.Portal>
  );
};

Popover.Close = PopoverPrimitive.Close;

