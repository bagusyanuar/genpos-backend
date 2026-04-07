CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES categories(id),
    level INTEGER NOT NULL DEFAULT 0,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100),
    description TEXT,
    type VARCHAR(20) NOT NULL DEFAULT 'PRODUCT',
    image_url TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_type_is_active ON categories(type, is_active);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id_sort_order ON categories(parent_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
