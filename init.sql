Create Table IF NOT EXISTS commands(
    id serial PRIMARY KEY,
    command TEXT,
    result TEXT,
    status VARCHAR(255)
);
