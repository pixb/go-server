import React, { useState } from 'react';

interface LoginFormProps {
  onSubmit: (username: string, password: string) => Promise<void>;
  loading?: boolean;
  error?: string;
}

export const LoginForm: React.FC<LoginFormProps> = ({ onSubmit, loading = false, error }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [validationError, setValidationError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setValidationError('');

    if (!username.trim()) {
      setValidationError('用户名不能为空');
      return;
    }

    if (!password.trim()) {
      setValidationError('密码不能为空');
      return;
    }

    if (password.length < 6) {
      setValidationError('密码长度至少6位');
      return;
    }

    await onSubmit(username, password);
  };

  return (
    <form onSubmit={handleSubmit} className="login-form" data-testid="login-form">
      <h2 data-testid="login-title">登录</h2>
      
      {(error || validationError) && (
        <div className="error-message" data-testid="error-message" role="alert">
          {error || validationError}
        </div>
      )}

      <div className="form-group">
        <label htmlFor="username">用户名</label>
        <input
          id="username"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          disabled={loading}
          data-testid="username-input"
          placeholder="请输入用户名"
        />
      </div>

      <div className="form-group">
        <label htmlFor="password">密码</label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={loading}
          data-testid="password-input"
          placeholder="请输入密码"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        data-testid="submit-button"
        className="submit-button"
      >
        {loading ? '登录中...' : '登录'}
      </button>
    </form>
  );
};

export default LoginForm;
