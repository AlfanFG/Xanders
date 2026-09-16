-- +goose Up
-- +goose StatementBegin
ALTER TABLE prompt_type ADD COLUMN icon TEXT DEFAULT '';
ALTER TABLE prompt_type ADD COLUMN background_class TEXT DEFAULT '';

INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Fashion', 'template', 'Clothing brands, fashion shows', '👗', ''),
('Accessories', 'template', 'Watches, bags, jewelry', '💎', ''),
('PC & Gaming', 'template', 'Custom rigs, peripherals', '🖥️', ''),
('Electronics', 'template', 'Smartphones, audio', '📱', ''),
('Automotive', 'template', 'Sports cars, luxury vehicles', '🏎️', '');

INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Golden Hour', 'lighting', '', '🌅', 'from-amber-500 to-orange-600'),
('Cyberpunk', 'lighting', '', '🏙️', 'from-pink-500 to-purple-600'),
('Soft Studio', 'lighting', '', '💡', 'from-slate-200 to-slate-400'),
('Dramatic', 'lighting', '', '🌗', 'from-gray-800 to-black');

INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Slow Tracking', 'camera', 'Smooth horizontal movement', '🛤️', ''),
('Drone Flyover', 'camera', 'High-angle cinematic sweep', '🚁', ''),
('Static Cinematic', 'camera', 'Locked-off tripod shot', '🎥', ''),
('Orbit 360', 'camera', 'Circular panning around subject', '🔄', '');

INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Epic Orchestral', 'audio', '', '🎻', ''),
('Ambient Chill', 'audio', '', '🎹', ''),
('Corporate Upbeat', 'audio', '', '📈', ''),
('No Audio', 'audio', '', '🔇', '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM prompt_type WHERE type IN ('template', 'lighting', 'camera', 'audio');

ALTER TABLE prompt_type DROP COLUMN background_class;
ALTER TABLE prompt_type DROP COLUMN icon;
-- +goose StatementEnd
