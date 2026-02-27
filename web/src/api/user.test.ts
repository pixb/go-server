import { describe, test, expect, vi, beforeEach } from 'vitest';
import { registerUser } from './user';
import { Timestamp } from '@bufbuild/protobuf';

// Mock the connect client
vi.mock('@connectrpc/connect-web', () => ({
  createConnectTransport: vi.fn(() => ({})),
}));

vi.mock('@connectrpc/connect', () => ({
  createPromiseClient: vi.fn(() => ({
    registerUser: vi.fn(),
  })),
}));

describe('User API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('should register user successfully', async () => {
    const mockResponse = {
      user: {
        id: 1n,
        username: 'testuser',
        email: 'test@example.com',
        nickname: 'Test User',
        phone: '13800138000',
        role: 'user',
        createdAt: Timestamp.fromDate(new Date('2024-01-15T10:30:00Z')),
        updatedAt: Timestamp.fromDate(new Date('2024-01-15T10:30:00Z')),
      },
      accessToken: 'mock-access-token',
      refreshToken: 'mock-refresh-token',
      accessTokenExpiresAt: Timestamp.fromDate(new Date('2024-01-16T10:30:00Z')),
    };

    const { createPromiseClient } = await import('@connectrpc/connect');
    const mockClient = {
      registerUser: vi.fn().mockResolvedValue(mockResponse),
    };
    vi.mocked(createPromiseClient).mockReturnValue(mockClient as any);

    const result = await registerUser(
      'testuser',
      'test@example.com',
      'password123',
      'Test User',
      '13800138000'
    );

    expect(mockClient.registerUser).toHaveBeenCalledTimes(1);
    expect(result.user.username).toBe('testuser');
    expect(result.accessToken).toBe('mock-access-token');
    expect(result.user.createdAt).toBeDefined();
  });

  test('should handle registration error', async () => {
    const { createPromiseClient } = await import('@connectrpc/connect');
    const mockClient = {
      registerUser: vi.fn().mockRejectedValue(new Error('Username already exists')),
    };
    vi.mocked(createPromiseClient).mockReturnValue(mockClient as any);

    await expect(
      registerUser(
        'existinguser',
        'test@example.com',
        'password123',
        'Test User',
        '13800138000'
      )
    ).rejects.toThrow('Username already exists');
  });
});
