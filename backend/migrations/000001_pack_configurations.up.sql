CREATE TABLE IF NOT EXISTS pack_configurations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    pack_sizes INTEGER[] NOT NULL,
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default configuration with the main edge case
INSERT INTO pack_configurations (name, pack_sizes, is_default, is_active) 
VALUES ('Main Edge Case', ARRAY[23, 31, 53], true, true);

-- Create index for better query performance
CREATE INDEX idx_pack_configurations_is_default ON pack_configurations (is_default) WHERE is_default = true;
CREATE INDEX idx_pack_configurations_is_active ON pack_configurations (is_active) WHERE is_active = true; 