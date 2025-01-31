-- +goose Up
CREATE TABLE transactions
(
    id              SERIAL PRIMARY KEY,
    sender_id       INT            NOT NULL,
    receiver_id     INT            NOT NULL,
    amount          DECIMAL(15, 2) NOT NULL,
    transaction_type VARCHAR(50)   NOT NULL,
    created_at      TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;
