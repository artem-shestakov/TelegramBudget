CREATE TABLE IF NOT EXISTS budgets (
    id bigint PRIMARY KEY,
    title varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS incomes (
    id serial PRIMARY KEY,
    title varchar(255),
    plan numeric,
    budget_id bigint REFERENCES budgets(id) ON DELETE CASCADE
);