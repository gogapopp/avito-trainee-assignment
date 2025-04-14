CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE pvz (
    id TEXT PRIMARY KEY,
    registration_date TIMESTAMP WITH TIME ZONE NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE reception (
    id TEXT PRIMARY KEY,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    pvz_id TEXT NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE product (
    id TEXT PRIMARY KEY,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    type TEXT NOT NULL,
    reception_id TEXT NOT NULL REFERENCES reception(id) ON DELETE CASCADE
);

CREATE INDEX reception_pvz_id_idx ON reception(pvz_id);
CREATE INDEX reception_status_idx ON reception(status);
CREATE INDEX product_reception_id_idx ON product(reception_id);
