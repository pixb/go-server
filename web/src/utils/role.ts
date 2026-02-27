import { Role } from '../types/proto/api/v1/common_pb';

export function roleToString(role: Role): string {
  switch (role) {
    case Role.ADMIN:
      return '管理员';
    case Role.USER:
      return '用户';
    case Role.UNSPECIFIED:
    default:
      return '未知';
  }
}

export function isValidRole(role: Role): boolean {
  return role === Role.ADMIN || role === Role.USER;
}
