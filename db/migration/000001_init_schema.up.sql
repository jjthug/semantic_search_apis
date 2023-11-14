CREATE TABLE "users" (
                         "user_id" bigserial PRIMARY KEY,
                         "username" varchar UNIQUE NOT NULL,
                         "hashed_password" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "docs" (
                        "user_id" bigserial UNIQUE PRIMARY KEY NOT NULL,
                        "doc" varchar NOT NULL
);

ALTER TABLE "docs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");