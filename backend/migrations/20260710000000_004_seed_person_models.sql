-- +goose Up
-- +goose StatementBegin

-- Person Model archetypes for advertising videos
INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Professional Woman', 'person_model', 'Elegant, confident female model in smart professional attire', '👩‍💼', 'from-rose-400 to-pink-600'),
('Professional Man',   'person_model', 'Polished male model in business-casual look', '👨‍💼', 'from-blue-400 to-blue-700'),
('Young Athlete',      'person_model', 'Energetic, fit model in sportswear — conveys performance and motion', '🏃', 'from-orange-400 to-red-600'),
('Creative Artist',    'person_model', 'Edgy, expressive model with artistic flair and bold styling', '🎨', 'from-purple-400 to-violet-700'),
('Lifestyle Influencer','person_model','Trendy, relatable everyday person in casual modern settings', '📸', 'from-teal-400 to-cyan-600'),
('Luxury Model',       'person_model', 'High-fashion editorial-style presence, haute couture aesthetic', '💎', 'from-yellow-400 to-amber-600'),
('Tech Enthusiast',    'person_model', 'Modern, casual persona interacting with smart devices and gadgets', '💻', 'from-sky-400 to-indigo-600'),
('No Model (Product Only)', 'person_model', 'Focus purely on the product with no person in frame', '📦', 'from-slate-600 to-slate-800');

-- Additional template categories (on top of existing: Fashion, Accessories, PC & Gaming, Electronics, Automotive)
INSERT INTO prompt_type (name, type, description, icon, background_class) VALUES
('Furniture',          'template', 'Sofas, tables, chairs, interior design pieces', '🛋️', ''),
('Home Goods',         'template', 'Kitchen appliances, home decor, household essentials', '🏠', ''),
('Food & Beverage',    'template', 'Restaurants, packaged food, beverages, culinary products', '🍽️', ''),
('Beauty & Skincare',  'template', 'Cosmetics, serums, skincare routines, makeup brands', '💄', ''),
('Health & Wellness',  'template', 'Supplements, gym equipment, yoga, self-care products', '🧘', ''),
('Sports & Outdoors',  'template', 'Sporting gear, camping equipment, adventure products', '⛺', ''),
('Footwear',           'template', 'Sneakers, boots, luxury shoes, sandals', '👟', ''),
('Luxury Goods',       'template', 'Premium watches, exclusive designer brands, rare collectibles', '⌚', ''),
('Travel & Hospitality','template','Hotels, airlines, travel destinations, resort experiences', '✈️', ''),
('Kids & Toys',        'template', 'Children''s products, educational toys, baby goods', '🧸', ''),
('Pet Products',       'template', 'Pet food, accessories, grooming and care items', '🐾', ''),
('Real Estate',        'template', 'Property listings, architecture, interior design showcases', '🏡', '');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM prompt_type WHERE type = 'person_model';
DELETE FROM prompt_type WHERE type = 'template' AND name IN (
  'Furniture', 'Home Goods', 'Food & Beverage', 'Beauty & Skincare',
  'Health & Wellness', 'Sports & Outdoors', 'Footwear', 'Luxury Goods',
  'Travel & Hospitality', 'Kids & Toys', 'Pet Products', 'Real Estate'
);
-- +goose StatementEnd
