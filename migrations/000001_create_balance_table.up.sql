CREATE TABLE IF NOT EXISTS balance(
    user_id INT PRIMARY KEY,
    balance decimal NOT NULL
    CHECK (balance>=0)
);