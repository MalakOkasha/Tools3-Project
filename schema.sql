CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
name VARCHAR(255) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
phone VARCHAR(50) NOT NULL,
location VARCHAR(255),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE couriers (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
vehicle_type VARCHAR(50) NOT NULL,
available BOOLEAN DEFAULT TRUE,
last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
orders UUID[] DEFAULT '{}',  -- Array of UUIDs to store order IDs
store_id UUID,  -- UUID for store identification
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE SET NULL  -- Link store_id to stores table
);



CREATE TABLE admins (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,  -- Foreign key to link to the users table
store_id UUID,  -- UUID for store identification
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE SET NULL  -- Link store_id to stores table
);



CREATE TABLE stores (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
name VARCHAR(255) NOT NULL,
location VARCHAR(255),
owner_id UUID,  -- Foreign key referencing the user who owns the store
admins_ids UUID[] DEFAULT '{}',  -- Array of UUIDs to store admin IDs
couriers_ids UUID[] DEFAULT '{}',  -- Array of UUIDs to store courier IDs
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL
);



CREATE TABLE owners (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,                   -- Foreign key linking to the users table
store_id UUID,                           -- Foreign key linking to the stores table
store_name VARCHAR(255) UNIQUE NOT NULL, -- Store name with a UNIQUE constraint
store_location VARCHAR(255),             -- Store location specific to the owner
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,      -- Cascade deletion if user is deleted
FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE SET NULL    -- Set to NULL if the store is deleted
);







CREATE TABLE items (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,                         -- The owner (seller) who created the item
store_id UUID NOT NULL,                        -- Foreign key linking to the stores table
name VARCHAR(255) NOT NULL,                    -- Item name
description TEXT,                              -- Detailed description of the item
price DECIMAL(10, 2) NOT NULL CHECK (price > 0), -- Price of the item, must be greater than 0
stock INT DEFAULT 0 CHECK (stock >= 0),        -- Quantity of the item in stock
category VARCHAR(100),                         -- Category of the item
cover_link VARCHAR(255),                       -- URL for the main cover image of the item
images TEXT[] DEFAULT '{}',                    -- Array of image URLs
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of item creation
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of last update
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);







CREATE TABLE orders (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
courier_id UUID,
store_id UUID NOT NULL,
item_ids UUID[] NOT NULL,
total_price NUMERIC(10, 2) NOT NULL,
status VARCHAR(20) NOT NULL DEFAULT 'pending',
pickup_location TEXT NOT NULL,
drop_off_location TEXT NOT NULL,
package_details TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
delivered_at TIMESTAMPTZ,
-- Foreign key constraints
CONSTRAINT fk_user
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
CONSTRAINT fk_courier
FOREIGN KEY (courier_id) REFERENCES couriers (id) ON DELETE SET NULL,
CONSTRAINT fk_store
FOREIGN KEY (store_id) REFERENCES stores (id) ON DELETE CASCADE
);