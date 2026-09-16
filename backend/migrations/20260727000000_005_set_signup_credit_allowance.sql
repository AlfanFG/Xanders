-- +goose Up
-- New accounts receive a one-time 15-credit signup allowance. Existing
-- balances are intentionally preserved: reducing a previously granted balance
-- requires an explicit product/admin decision.
ALTER TABLE users ALTER COLUMN credit_balance SET DEFAULT 15;

-- +goose Down
ALTER TABLE users ALTER COLUMN credit_balance SET DEFAULT 0;
