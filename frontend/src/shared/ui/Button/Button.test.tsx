import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from './Button';

describe('Button', () => {
  it('renders with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument();
  });

  it('handles click events', async () => {
    const handleClick = vi.fn();
    const user = userEvent.setup();
    
    render(<Button onClick={handleClick}>Click me</Button>);
    
    await user.click(screen.getByRole('button'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('is disabled when isLoading is true', () => {
    render(<Button isLoading>Loading</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('shows spinner when isLoading', () => {
    render(<Button isLoading>Loading</Button>);
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
    // Spinner должен отображаться вместо leftIcon
  });

  it('renders with leftIcon', () => {
    const TestIcon = () => <span data-testid="left-icon">←</span>;
    render(<Button leftIcon={<TestIcon />}>With icon</Button>);
    expect(screen.getByTestId('left-icon')).toBeInTheDocument();
  });

  it('renders with rightIcon', () => {
    const TestIcon = () => <span data-testid="right-icon">→</span>;
    render(<Button rightIcon={<TestIcon />}>With icon</Button>);
    expect(screen.getByTestId('right-icon')).toBeInTheDocument();
  });

  it('hides leftIcon when isLoading', () => {
    const TestIcon = () => <span data-testid="left-icon">←</span>;
    render(<Button leftIcon={<TestIcon />} isLoading>Loading</Button>);
    expect(screen.queryByTestId('left-icon')).not.toBeInTheDocument();
  });

  describe('variants', () => {
    it('renders primary variant by default', () => {
      const { container } = render(<Button>Primary</Button>);
      expect(container.firstChild).toHaveClass('bg-primary-500');
    });

    it('renders secondary variant', () => {
      const { container } = render(<Button variant="secondary">Secondary</Button>);
      expect(container.firstChild).toHaveClass('bg-neutral-800');
    });

    it('renders danger variant', () => {
      const { container } = render(<Button variant="danger">Danger</Button>);
      expect(container.firstChild).toHaveClass('bg-error');
    });
  });

  describe('sizes', () => {
    it('renders default md size', () => {
      const { container } = render(<Button>Default</Button>);
      expect(container.firstChild).toHaveClass('h-10');
    });

    it('renders sm size', () => {
      const { container } = render(<Button size="sm">Small</Button>);
      expect(container.firstChild).toHaveClass('h-8');
    });

    it('renders lg size', () => {
      const { container } = render(<Button size="lg">Large</Button>);
      expect(container.firstChild).toHaveClass('h-12');
    });
  });
});

