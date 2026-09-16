-- +goose Up
-- +goose StatementBegin
CREATE TABLE prompt_type (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    type        TEXT NOT NULL DEFAULT 'template',  -- template | audio | lighting | camera
    description TEXT NOT NULL DEFAULT '',
    is_deleted  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prompt_type;
-- +goose StatementEnd
