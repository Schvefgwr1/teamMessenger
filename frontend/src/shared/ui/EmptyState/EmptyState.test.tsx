import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { EmptyState } from './EmptyState';
import { Search, MessageSquare } from 'lucide-react';

describe('EmptyState', () => {
  it('renders with icon', () => {
    render(<EmptyState icon={Search} title="No results" />);
    expect(screen.getByText('No results')).toBeInTheDocument();
  });

  it('renders with title', () => {
    render(<EmptyState title="No data" />);
    expect(screen.getByText('No data')).toBeInTheDocument();
  });

  it('renders with description', () => {
    render(<EmptyState description="Try again later" />);
    expect(screen.getByText('Try again later')).toBeInTheDocument();
  });

  it('renders action button', async () => {
    const handleClick = vi.fn();
    const user = userEvent.setup();

    render(
      <EmptyState
        title="Empty"
        action={{
          label: 'Create',
          onClick: handleClick,
        }}
      />
    );

    const button = screen.getByRole('button', { name: /create/i });
    expect(button).toBeInTheDocument();

    await user.click(button);
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('renders action with icon', () => {
    render(
      <EmptyState
        action={{
          label: 'Create',
          onClick: () => {},
          icon: <Search data-testid="action-icon" />,
        }}
      />
    );

    expect(screen.getByTestId('action-icon')).toBeInTheDocument();
  });

  it('renders children', () => {
    render(
      <EmptyState>
        <div data-testid="custom-content">Custom content</div>
      </EmptyState>
    );

    expect(screen.getByTestId('custom-content')).toBeInTheDocument();
  });

  describe('sizes', () => {
    it('renders with sm size', () => {
      const { container } = render(
        <EmptyState size="sm" icon={Search} title="Small" />
      );
      expect(container.firstChild).toHaveClass('py-12');
    });

    it('renders with md size by default', () => {
      const { container } = render(
        <EmptyState icon={Search} title="Medium" />
      );
      expect(container.firstChild).toHaveClass('py-12');
    });

    it('renders with lg size', () => {
      const { container } = render(
        <EmptyState size="lg" icon={Search} title="Large" />
      );
      expect(container.firstChild).toHaveClass('py-12');
    });
  });

  it('applies custom className', () => {
    const { container } = render(
      <EmptyState className="custom-class" title="Test" />
    );
    expect(container.firstChild).toHaveClass('custom-class');
  });
});

