CREATE TABLE users (
    upn TEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    password TEXT,
    mfa_token TEXT,
    azuread_oid TEXT UNIQUE
);

CREATE TYPE device_state AS ENUM ('deploying', 'managed', 'user_unenrolled', 'missing');
CREATE TYPE enrollment_type AS ENUM ('Unenrolled', 'User', 'Device');

CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    udid TEXT UNIQUE NOT NULL,
    state device_state NOT NULL,
    enrollment_type enrollment_type NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    model TEXT DEFAULT '' NOT NULL,
    hw_dev_id TEXT UNIQUE NOT NULL,
    operating_system TEXT NOT NULL,
    azure_did TEXT UNIQUE,
    nodecache_version TEXT DEFAULT '' NOT NULL,
    lastseen TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    lastseen_status INTEGER DEFAULT 0 NOT NULL,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    enrolled_by TEXT REFERENCES users(upn)
);

CREATE TABLE device_inventory (
	id SERIAL PRIMARY KEY,
    device_id INTEGER REFERENCES devices(id) NOT NULL,
    uri TEXT NOT NULL,
    format TEXT DEFAULT '' NOT NULL,
    value TEXT DEFAULT '' NOT NULL,
    UNIQUE (device_id, uri)
);

CREATE TABLE device_session_cache (

);

CREATE TABLE policies (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    priority SMALLINT DEFAULT '0' NOT NULL
);

CREATE TABLE policies_payload (
    id SERIAL PRIMARY KEY,
    policy_id INTEGER REFERENCES policies(id),
    uri TEXT NOT NULL,
    format TEXT DEFAULT '' NOT NULL,
    type TEXT DEFAULT '' NOT NULL,
    value TEXT DEFAULT '' NOT NULL,
    exec BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (policy_id, uri)
);

CREATE TABLE device_cache (
    device_id INTEGER REFERENCES devices(id) NOT NULL,
    payload_id INTEGER REFERENCES policies_payload(id),
    inventory_id INTEGER REFERENCES device_inventory(id),
    cache_id SERIAL NOT NULL,
    PRIMARY KEY (device_id, cache_id),
    CONSTRAINT chk_reference check ((payload_id is not null and inventory_id is null) or (payload_id is null and inventory_id is not null))
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    priority SMALLINT DEFAULT '0' NOT NULL
);

CREATE TABLE group_devices (
    group_id INTEGER REFERENCES groups(id) NOT NULL,
    device_id INTEGER REFERENCES devices(id) NOT NULL,
    PRIMARY KEY (group_id, device_id)
);

CREATE TABLE group_policies (
    group_id INTEGER REFERENCES groups(id),
    policy_id INTEGER REFERENCES policies(id),
    PRIMARY KEY (group_id, policy_id)
);

CREATE TABLE settings (
    tenant_name TEXT NOT NULL,
    tenant_email TEXT DEFAULT '' NOT NULL,
    tenant_website TEXT DEFAULT '' NOT NULL,
    tenant_phone TEXT DEFAULT '' NOT NULL,
    tenant_azureid TEXT NOT NULL,
    disable_enrollment BOOLEAN DEFAULT false NOT NULL
);

CREATE TABLE certificates (
    id TEXT PRIMARY KEY,
    cert BYTEA,
    key BYTEA
);