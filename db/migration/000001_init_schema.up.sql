CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "passwordHash" varchar NOT NULL,
                         "phone" string NOT NULL,
                         "email" string NOT NULL,
                         "private_contact" boolean NOT NULL,
                         "about_id" integer UNIQUE
);

CREATE TABLE "abouts" (
                          "id" bigserial UNIQUE PRIMARY KEY,
                          "about_description" text NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("about_id") REFERENCES "abouts" ("id");