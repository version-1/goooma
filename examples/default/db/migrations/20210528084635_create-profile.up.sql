CREATE TABLE profiles (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id uuid NOT NULL,
  description text NOT NULL,
  avatar_id uuid,
  avatar json,
  location json NOT NULL,

  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_id_on_profiles ON profiles(user_id);
CREATE INDEX idx_avatar_id_on_profiles ON profiles(avatar_id);
