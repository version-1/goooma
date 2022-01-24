CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  display_name varchar NOT NULL,
  email        varchar UNIQUE,

  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_display_name_on_users ON users(display_name);
CREATE INDEX idx_email_on_users ON users(email);
