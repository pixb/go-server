import { describe, test, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginForm } from './LoginForm';

describe('LoginForm', () => {
  test('should render login form correctly', () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    expect(screen.getByTestId('login-title')).toHaveTextContent('登录');
    expect(screen.getByTestId('username-input')).toBeInTheDocument();
    expect(screen.getByTestId('password-input')).toBeInTheDocument();
    expect(screen.getByTestId('submit-button')).toHaveTextContent('登录');
  });

  test('should show validation error when username is empty', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    const submitButton = screen.getByTestId('submit-button');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('用户名不能为空');
    });
  });

  test('should show validation error when password is empty', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    const usernameInput = screen.getByTestId('username-input');
    fireEvent.change(usernameInput, { target: { value: 'testuser' } });

    const submitButton = screen.getByTestId('submit-button');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('密码不能为空');
    });
  });

  test('should show validation error when password is too short', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    const usernameInput = screen.getByTestId('username-input');
    const passwordInput = screen.getByTestId('password-input');

    fireEvent.change(usernameInput, { target: { value: 'testuser' } });
    fireEvent.change(passwordInput, { target: { value: '123' } });

    const submitButton = screen.getByTestId('submit-button');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('密码长度至少6位');
    });
  });

  test('should call onSubmit with correct credentials', async () => {
    const mockSubmit = vi.fn().mockResolvedValue(undefined);
    render(<LoginForm onSubmit={mockSubmit} />);

    const usernameInput = screen.getByTestId('username-input');
    const passwordInput = screen.getByTestId('password-input');

    fireEvent.change(usernameInput, { target: { value: 'testuser' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });

    const submitButton = screen.getByTestId('submit-button');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalledWith('testuser', 'password123');
    });
  });

  test('should show loading state', () => {
    render(<LoginForm onSubmit={vi.fn()} loading={true} />);

    const submitButton = screen.getByTestId('submit-button');
    expect(submitButton).toHaveTextContent('登录中...');
    expect(submitButton).toBeDisabled();

    const usernameInput = screen.getByTestId('username-input');
    const passwordInput = screen.getByTestId('password-input');
    expect(usernameInput).toBeDisabled();
    expect(passwordInput).toBeDisabled();
  });

  test('should display error message from props', () => {
    render(<LoginForm onSubmit={vi.fn()} error="登录失败，请重试" />);

    expect(screen.getByTestId('error-message')).toHaveTextContent('登录失败，请重试');
  });

  test('should clear validation error when user starts typing', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    // Trigger validation error
    const submitButton = screen.getByTestId('submit-button');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toBeInTheDocument();
    });

    // Type in username field
    const usernameInput = screen.getByTestId('username-input');
    fireEvent.change(usernameInput, { target: { value: 'testuser' } });

    // Submit again to trigger new validation
    fireEvent.click(submitButton);

    // The error should now be about password
    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('密码不能为空');
    });
  });

  test('should have correct input types', () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    const usernameInput = screen.getByTestId('username-input');
    const passwordInput = screen.getByTestId('password-input');

    expect(usernameInput).toHaveAttribute('type', 'text');
    expect(passwordInput).toHaveAttribute('type', 'password');
  });

  test('should have correct placeholders', () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    const usernameInput = screen.getByTestId('username-input');
    const passwordInput = screen.getByTestId('password-input');

    expect(usernameInput).toHaveAttribute('placeholder', '请输入用户名');
    expect(passwordInput).toHaveAttribute('placeholder', '请输入密码');
  });
});
