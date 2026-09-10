USE game_store;

-- Data Reset
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE order_details;
TRUNCATE TABLE orders;
TRUNCATE TABLE carts;
TRUNCATE TABLE game_keys;
TRUNCATE TABLE games;
TRUNCATE TABLE categories;
TRUNCATE TABLE user_profiles;
TRUNCATE TABLE users;
SET FOREIGN_KEY_CHECKS = 1;

-- Seed Users 
-- Admin: admin@store.com | Password: admin123
-- Customer: buyer@gmail.com | Password: password123
INSERT INTO users (id, email, password, role) VALUES
(1, 'admin@store.com', 'admin123', 'admin'),
(2, 'buyer@gmail.com', 'password123', 'customer');

-- Seed User Profiles
INSERT INTO user_profiles (user_id, full_name, wallet_balance) VALUES
(1, 'Admin Game Store', 0.00),
(2, 'Budi Customer', 500000.00);

-- Seed Categories
INSERT INTO categories (id, category_name) VALUES
(1, 'Action'),
(2, 'RPG'),
(3, 'Strategy');

-- Seed Games
INSERT INTO games (id, category_id, title, price) VALUES
(1, 1, 'Cyberpunk 2077', 350000.00),
(2, 2, 'Elden Ring', 450000.00),
(3, 3, 'Civilization VI', 150000.00);

-- Seed Game Keys (Lisensi Game)
INSERT INTO game_keys (game_id, license_key, status) VALUES
-- Key Cyberpunk 2077
(1, 'CP77-AAAA-1111', 'available'),
(1, 'CP77-BBBB-2222', 'available'),
-- Key Elden Ring
(2, 'ER33-CCCC-3333', 'available'),
(2, 'ER33-DDDD-4444', 'available'),
-- Key Civilization VI
(3, 'CIV6-EEEE-5555', 'available');