CREATE TABLE IF NOT EXISTS "addresses" (
    "id" SERIAL PRIMARY KEY,
    "hospital_id" INTEGER REFERENCES hospital(id),
    "region" VARCHAR(30),
    "street" VARCHAR(30)
);