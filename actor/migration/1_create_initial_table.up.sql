CREATE TABLE users (
  id uuid NOT NULL,
  first_name varchar(100) NOT NULL,
  middle_name varchar(100),
  last_name varchar(100) NOT NULL,
  phone varchar(50) NOT NULL,
  gender varchar(10) NOT NULL,
  birthdate date NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  timestamp timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (email)
);

CREATE TABLE address (
  id uuid NOT NULL,
  user_id uuid NOT NULL,
  street text NOT NULL,
  city varchar(100) NOT NULL,
  state varchar(100) NOT NULL,
  postal_code varchar(10) NOT NULL,
  country varchar(50) NOT NULL,
  timestamp timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (user_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE cascade ON UPDATE cascade
);

CREATE TABLE oauth (
  id uuid NOT NULL,
  consumer_id uuid NOT NULL,
  name varchar(20) NOT NULL,
  client_id varchar(255) NOT NULL,
  client_secret varchar(255) NOT NULL,
  timestamp timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

ALTER TABLE oauth ADD constraint oauth_consumer_id_users_id_foreign FOREIGN KEY (consumer_id) REFERENCES users(id)
ON DELETE cascade ON UPDATE cascade;