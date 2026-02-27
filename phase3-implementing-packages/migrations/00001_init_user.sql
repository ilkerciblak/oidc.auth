-- +goose Up 
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS users( 
  id UUID PRIMARY KEY,
	display_name VARCHAR(255) UNIQUE,
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(255) UNIQUE,
  last_login_at TIMESTAMPTZ,
  verified BOOLEAN
);

CREATE TABLE IF NOT EXISTS auth_provider(
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  provider VARCHAR(255) NOT NULL,
  provider_user_id VARCHAR(255) NOT NULL, 
  last_login_with TIMESTAMPTZ,
  
  CONSTRAINT fk_users
    FOREIGN KEY(user_id) 
    REFERENCES users(id)
    ON DELETE CASCADE,


  CONSTRAINT uq_provider_provider_id UNIQUE (provider, provider_user_id)
);


CREATE INDEX idx_user_providers ON auth_provider(provider, provider_user_id)



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SELECT 'down SQL query';
DROP INDEX idx_user_providers;
DROP TABLE IF EXISTS auth_provider;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
