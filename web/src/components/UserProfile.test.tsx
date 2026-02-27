import { describe, test, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { UserProfile } from './UserProfile';
import { User } from '../types/proto/api/v1/common_pb';
import { Timestamp } from '@bufbuild/protobuf';

describe('UserProfile', () => {
  const mockUser: User = {
    id: 1n,
    username: 'testuser',
    email: 'test@example.com',
    nickname: 'Test User',
    phone: '13800138000',
    role: 'user',
    createdAt: Timestamp.fromDate(new Date('2024-01-15T10:30:00Z')),
    updatedAt: Timestamp.fromDate(new Date('2024-02-20T14:45:00Z')),
    passwordExpiresAt: Timestamp.fromDate(new Date('2024-04-15T10:30:00Z')),
  };

  test('should render user profile correctly', () => {
    render(<UserProfile user={mockUser} />);

    // Check if user information is displayed
    expect(screen.getByTestId('user-nickname')).toHaveTextContent('Test User');
    expect(screen.getByTestId('user-username')).toHaveTextContent('testuser');
    expect(screen.getByTestId('user-email')).toHaveTextContent('test@example.com');
    expect(screen.getByTestId('user-phone')).toHaveTextContent('13800138000');
    expect(screen.getByTestId('user-role')).toHaveTextContent('user');
  });

  test('should display formatted timestamps', () => {
    render(<UserProfile user={mockUser} />);

    // Check if timestamps are formatted and displayed
    const createdAt = screen.getByTestId('user-created-at');
    const updatedAt = screen.getByTestId('user-updated-at');
    const passwordExpires = screen.getByTestId('user-password-expires');

    expect(createdAt).toBeInTheDocument();
    expect(updatedAt).toBeInTheDocument();
    expect(passwordExpires).toBeInTheDocument();

    // Check that the timestamps contain formatted date strings
    expect(createdAt.textContent).toContain('2024');
    expect(updatedAt.textContent).toContain('2024');
    expect(passwordExpires.textContent).toContain('2024');
  });

  test('should not show edit button when onUpdate is not provided', () => {
    render(<UserProfile user={mockUser} />);

    const editButton = screen.queryByTestId('edit-button');
    expect(editButton).not.toBeInTheDocument();
  });

  test('should show edit button when onUpdate is provided', () => {
    const mockOnUpdate = vi.fn();
    render(<UserProfile user={mockUser} onUpdate={mockOnUpdate} />);

    const editButton = screen.getByTestId('edit-button');
    expect(editButton).toBeInTheDocument();
    expect(editButton).toHaveTextContent('编辑');
  });

  test('should call onUpdate when edit button is clicked', () => {
    const mockOnUpdate = vi.fn();
    render(<UserProfile user={mockUser} onUpdate={mockOnUpdate} />);

    const editButton = screen.getByTestId('edit-button');
    fireEvent.click(editButton);

    expect(mockOnUpdate).toHaveBeenCalledTimes(1);
    expect(mockOnUpdate).toHaveBeenCalledWith(mockUser);
  });

  test('should handle user without password expiration', () => {
    const userWithoutPasswordExpiry: User = {
      ...mockUser,
      passwordExpiresAt: undefined,
    };

    render(<UserProfile user={userWithoutPasswordExpiry} />);

    const passwordExpires = screen.queryByTestId('user-password-expires');
    expect(passwordExpires).not.toBeInTheDocument();
  });

  test('should render with correct CSS classes', () => {
    render(<UserProfile user={mockUser} />);

    const profileContainer = screen.getByTestId('user-profile');
    expect(profileContainer).toHaveClass('user-profile');

    const passwordExpires = screen.getByTestId('user-password-expires');
    expect(passwordExpires).toHaveClass('password-expires');
  });
});
