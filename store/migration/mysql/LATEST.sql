-- system_setting
CREATE TABLE system_setting (
  name varchar(255) NOT NULL,
  value text NOT NULL,
  description text NOT NULL DEFAULT '',
  PRIMARY KEY (name)
);

-- users table
CREATE TABLE users (
  id BIGINT AUTO_INCREMENT NOT NULL,
  username varchar(50) NOT NULL,
  nickname varchar(50) NULL,
  `password` varchar(255) NOT NULL,
  phone varchar(20) NULL,
  email varchar(100) NULL,
  `role` varchar(20) DEFAULT 'user',
  password_expires DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_email (email),
  UNIQUE KEY idx_users_username (username),
  KEY idx_users_deleted_at (deleted_at)
);

-- refresh_tokens table
CREATE TABLE refresh_tokens (
  id BIGINT AUTO_INCREMENT NOT NULL,
  user_id bigint NOT NULL,
  token text NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked boolean DEFAULT false,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_refresh_tokens_token (token(255)),
  KEY idx_refresh_tokens_deleted_at (deleted_at),
  KEY idx_refresh_tokens_user_id (user_id),
  CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
