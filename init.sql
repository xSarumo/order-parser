CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
    transaction_id VARCHAR(50) PRIMARY KEY,
    request_id VARCHAR(50),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id BIGINT PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(10) NOT NULL,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    payment_transaction_id VARCHAR(50) REFERENCES payments(transaction_id) ON DELETE CASCADE,
    locale VARCHAR(5) NOT NULL,
    internal_signature VARCHAR(50),
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT REFERENCES items(chrt_id) ON DELETE CASCADE,
    PRIMARY KEY (order_uid, chrt_id)
);