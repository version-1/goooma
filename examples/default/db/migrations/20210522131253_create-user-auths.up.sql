CREATE TABLE user_auths (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id uuid NOT NULL,
  provider varchar NOT NULL,
  uid varchar NOT NULL,
  access_token varchar NOT NULL,
  refresh_token varchar NOT NULL,

  expired_at timestamp NOT NULL,

  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW(),

  UNIQUE(provider, uid)
);

CREATE INDEX idx_user_id_on_user_auths ON user_auths(user_id);
CREATE INDEX idx_uid_on_user_auths ON user_auths(uid);
