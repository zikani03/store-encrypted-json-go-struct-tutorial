CREATE TABLE settings (
    id serial primary key not null,
    key text not null,
    data_json json not null,
    last_updated_at timestamp
);
