CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    price DOUBLE NOT NULL,
    tax DOUBLE NOT NULL,
    final_price DOUBLE NOT NULL
);