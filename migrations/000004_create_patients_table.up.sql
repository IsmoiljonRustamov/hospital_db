CREATE TABLE IF NOT EXISTS "patients" (
    "id" SERIAL PRIMARY KEY,
    "hospital_id" INTEGER REFERENCES hospital(id),
    "full_name" VARCHAR(40) NOT NULL,
    "patient_info" TEXT,
    "phone_number" TEXT
);