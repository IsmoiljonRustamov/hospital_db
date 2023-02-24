CREATE TABLE IF NOT EXISTS "staff" (
    "id" SERIAL PRIMARY KEY ,
    "hospital_id" INTEGER REFERENCES hospital(id),
    "full_name" TEXT,
    "phone_numbers" TEXT
);