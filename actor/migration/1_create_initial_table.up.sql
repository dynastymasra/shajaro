CREATE TYPE gender AS ENUM ('Male', 'Female');

CREATE TABLE users (
  id uuid NOT NULL,
  consumer_id uuid NOT NULL,
  first_name varchar(100) NOT NULL,
  middle_name varchar(100),
  last_name varchar(100) NOT NULL,
  avatar_url varchar(255) NOT NULL,
  phone jsonb NOT NULL,
  gender gender NOT NULL,
  address jsonb NOT NULL,
  birth_date date NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  timestamp timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (email, consumer_id)
);

-- CREATE TYPE oauth_name AS ENUM ('Android', 'iOS', 'Web', 'Desktop');
--
-- CREATE TABLE consumers (
--   id uuid NOT NULL,
--   consumer_id uuid NOT NULL,
--   user_id uuid NOT NULL,
--   name oauth_name NOT NULL,
--   client_id varchar(255) NOT NULL,
--   client_secret varchar(255) NOT NULL,
--   timestamp timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   PRIMARY KEY (id),
--   FOREIGN KEY (user_id) REFERENCES users(id)
-- );
--
-- ALTER TABLE consumers ADD constraint consumers_user_id_users_id_foreign FOREIGN KEY (user_id) REFERENCES users(id)
-- ON DELETE cascade ON UPDATE cascade;