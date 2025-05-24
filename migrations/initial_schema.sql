CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    chat_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TYPE config_status AS ENUM ('new', 'paid', 'expired', 'deletion');

CREATE TABLE IF NOT EXISTS configurations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status config_status DEFAULT 'new',
    expiration_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    interface_id INTEGER NOT NULL REFERENCES servers(id),
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    allowed_ip VARCHAR(255) NOT NULL,
    latest_handshake TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id)
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_configurations_user_id ON configurations(user_id);
CREATE INDEX idx_configurations_interface_id ON configurations(interface_id);
CREATE INDEX idx_configurations_status ON configurations(status);
CREATE INDEX idx_configurations_expiration_time ON configurations(expiration_time);
