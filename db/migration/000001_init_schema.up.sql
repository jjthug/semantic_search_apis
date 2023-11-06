CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "name" varchar NOT NULL,
                         "password_hash" varchar NOT NULL,
                         "phone" varchar NOT NULL,
                         "email" varchar NOT NULL,
                         "private_contact" boolean NOT NULL,
                         "about_description" text NOT NULL
);