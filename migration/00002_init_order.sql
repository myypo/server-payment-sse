-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM (
    'cool_order_created',
    'confirmed_by_mayor',
    'sbu_verification_pending',
    'changed_my_mind',
    'failed',
    'give_my_money_back',
    'chinazes'
);

CREATE TABLE "order" (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL
);


CREATE TABLE "event_order" (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID REFERENCES "order" (id) ON DELETE CASCADE,
    status ORDER_STATUS NOT NULL,

    created_at TIMESTAMPTZ NOT NULL
);

CREATE OR REPLACE FUNCTION is_order_final(
    order_status ORDER_STATUS,
    order_updated_at TIMESTAMPTZ,
    confirm_in INTERVAL
)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN (
        order_status IN ('changed_my_mind', 'failed', 'give_my_money_back')
        OR (order_status = 'chinazes' AND order_updated_at < now() - confirm_in)
    );
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
DROP TYPE IF EXISTS ORDER_STATUS;
DROP TABLE IF EXISTS "event_order";
DROP TYPE IF EXISTS EVENT_ORDER_STATUS;
DROP FUNCTION IF EXISTS is_order_final (VARCHAR(256), TIMESTAMPTZ, INTERVAL);
-- +goose StatementEnd
