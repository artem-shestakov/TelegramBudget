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

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    amount INT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS top_ups (
    income_id int REFERENCES incomes(id) ON DELETE CASCADE
) INHERITS (transactions);