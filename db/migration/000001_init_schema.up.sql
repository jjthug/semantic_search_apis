CREATE TABLE "users" (
                         "user_id" bigserial UNIQUE PRIMARY KEY NOT NULL,
                         "username" varchar UNIQUE NOT NULL,
                         "hashed_password" varchar NOT NULL,
                         "full_name" varchar NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "password_changed_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "docs" (
                        "user_id" bigserial UNIQUE PRIMARY KEY NOT NULL,
                        "doc" varchar NOT NULL
);

ALTER TABLE "docs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");