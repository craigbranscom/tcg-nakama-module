CREATE TABLE IF NOT EXISTS players (
    player_id UUID PRIMARY KEY,
    display_name TEXT,
    inventory JSONB
)

INSERT INTO players (player_id, display_name, inventory)
VALUES (
    'test_player',
    'test_display_name',
    '{"player_id": "test_player", "inventory": []}'
)