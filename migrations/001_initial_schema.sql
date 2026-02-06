-- Migration: 001_initial_schema
-- Description: Create initial database schema for Asset Manager

-- Users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
);

-- Asset types table
CREATE TABLE IF NOT EXISTS asset_types (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_asset_types_name (name),
    INDEX idx_asset_types_deleted_at (deleted_at)
);

-- Assets table
CREATE TABLE IF NOT EXISTS assets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    asset_type_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(255),
    serial_number VARCHAR(255),
    order_no VARCHAR(100),
    license_number VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (asset_type_id) REFERENCES asset_types(id),
    INDEX idx_assets_asset_type_id (asset_type_id),
    INDEX idx_assets_name (name),
    INDEX idx_assets_serial_number (serial_number),
    INDEX idx_assets_deleted_at (deleted_at)
);

-- Properties table (custom properties for assets)
CREATE TABLE IF NOT EXISTS properties (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    data_type ENUM('string', 'int', 'decimal', 'boolean', 'date', 'datetime', 'enum') NOT NULL DEFAULT 'string',
    enum_options JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_properties_name (name),
    INDEX idx_properties_deleted_at (deleted_at)
);

-- Assets properties junction table
CREATE TABLE IF NOT EXISTS assets_properties (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    property_id BIGINT NOT NULL,
    value TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id),
    FOREIGN KEY (property_id) REFERENCES properties(id),
    INDEX idx_assets_properties_asset_id (asset_id),
    INDEX idx_assets_properties_property_id (property_id),
    INDEX idx_assets_properties_deleted_at (deleted_at)
);

-- Persons table
CREATE TABLE IF NOT EXISTS persons (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_persons_name (name),
    INDEX idx_persons_email (email),
    INDEX idx_persons_deleted_at (deleted_at)
);

-- Attributes table (custom attributes for persons)
CREATE TABLE IF NOT EXISTS attributes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    data_type ENUM('string', 'int', 'decimal', 'boolean', 'date', 'datetime', 'enum') NOT NULL DEFAULT 'string',
    enum_options JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_attributes_name (name),
    INDEX idx_attributes_deleted_at (deleted_at)
);

-- Persons attributes junction table
CREATE TABLE IF NOT EXISTS persons_attributes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    person_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    value TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (person_id) REFERENCES persons(id),
    FOREIGN KEY (attribute_id) REFERENCES attributes(id),
    INDEX idx_persons_attributes_person_id (person_id),
    INDEX idx_persons_attributes_attribute_id (attribute_id),
    INDEX idx_persons_attributes_deleted_at (deleted_at)
);

-- Asset assignments table
CREATE TABLE IF NOT EXISTS asset_assignments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    person_id BIGINT NOT NULL,
    effective_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id),
    FOREIGN KEY (person_id) REFERENCES persons(id),
    INDEX idx_asset_assignments_asset_id (asset_id),
    INDEX idx_asset_assignments_person_id (person_id),
    INDEX idx_asset_assignments_effective_from (effective_from),
    INDEX idx_asset_assignments_effective_to (effective_to),
    INDEX idx_asset_assignments_deleted_at (deleted_at)
);

-- Insert default admin user (password: admin)
-- Password hash is bcrypt of 'admin' - generated with cost 10
-- IMPORTANT: Change this password immediately after first login!
INSERT INTO users (username, email, password_hash, is_active)
VALUES ('admin', 'assets@localhost', '$2a$10$rDkPvvAFV8kqwvKJzwlJAOHYnAHmFT4dNp.VKXBK5xgplqKje8.Hy', TRUE);

-- Insert special 'Unassigned' person for stock management
INSERT INTO persons (name, email, phone)
VALUES ('Unassigned', '', '');
