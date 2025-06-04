CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     first_name VARCHAR(50) NOT NULL,
                                     last_name VARCHAR(50) NOT NULL,
                                     login VARCHAR(50) UNIQUE NOT NULL,
                                     full_name VARCHAR(100),
                                     age INT NOT NULL,
                                     is_married BOOLEAN DEFAULT false,
                                     password VARCHAR(100) NOT NULL,
                                     role VARCHAR(20) DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        description TEXT NOT NULL,
                                        tags TEXT[],
                                        quantity INT NOT NULL,
                                        price NUMERIC(10,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INT NOT NULL REFERENCES users(id),
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
                                           id SERIAL PRIMARY KEY,
                                           order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                                           product_id INT NOT NULL REFERENCES products(id),
                                           quantity INT NOT NULL,
                                           price_at_order NUMERIC(10,2) NOT NULL
);