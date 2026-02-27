import { describe, test, expect } from 'vitest';
import { Role } from '../types/proto/api/v1/common_pb';
import { roleToString, isValidRole } from './role';

describe('Role Utilities', () => {
  describe('roleToString', () => {
    test('should convert ADMIN role to Chinese', () => {
      expect(roleToString(Role.ADMIN)).toBe('管理员');
    });

    test('should convert USER role to Chinese', () => {
      expect(roleToString(Role.USER)).toBe('用户');
    });

    test('should convert UNSPECIFIED role to Chinese', () => {
      expect(roleToString(Role.UNSPECIFIED)).toBe('未知');
    });
  });

  describe('isValidRole', () => {
    test('should return true for ADMIN role', () => {
      expect(isValidRole(Role.ADMIN)).toBe(true);
    });

    test('should return true for USER role', () => {
      expect(isValidRole(Role.USER)).toBe(true);
    });

    test('should return false for UNSPECIFIED role', () => {
      expect(isValidRole(Role.UNSPECIFIED)).toBe(false);
    });
  });
});
