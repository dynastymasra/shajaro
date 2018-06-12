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