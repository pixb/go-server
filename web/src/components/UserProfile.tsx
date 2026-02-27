import React from 'react';
import { User } from '../types/proto/api/v1/common_pb';
import { formatTimestamp } from '../utils/time';
import { roleToString } from '../utils/role';

interface UserProfileProps {
  user: User;
  onUpdate?: (user: User) => void;
}

export const UserProfile: React.FC<UserProfileProps> = ({ user, onUpdate }) => {
  return (
    <div className="user-profile" data-testid="user-profile">
      <h2 data-testid="user-nickname">{user.nickname}</h2>
      <div className="user-info">
        <p data-testid="user-username">
          <strong>用户名:</strong> {user.username}
        </p>
        <p data-testid="user-email">
          <strong>邮箱:</strong> {user.email}
        </p>
        <p data-testid="user-phone">
          <strong>电话:</strong> {user.phone}
        </p>
        <p data-testid="user-role">
          <strong>角色:</strong> {roleToString(user.role)}
        </p>
        <p data-testid="user-created-at">
          <strong>创建时间:</strong> {formatTimestamp(user.createdAt)}
        </p>
        <p data-testid="user-updated-at">
          <strong>更新时间:</strong> {formatTimestamp(user.updatedAt)}
        </p>
        {user.passwordExpiresAt && (
          <p data-testid="user-password-expires" className="password-expires">
            <strong>密码过期时间:</strong> {formatTimestamp(user.passwordExpiresAt)}
          </p>
        )}
      </div>
      {onUpdate && (
        <button
          data-testid="edit-button"
          onClick={() => onUpdate(user)}
          className="edit-button"
        >
          编辑
        </button>
      )}
    </div>
  );
};

export default UserProfile;
