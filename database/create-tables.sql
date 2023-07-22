DROP TABLE IF EXISTS AuthInfo;
DROP TABLE IF EXISTS UserInfo;

CREATE TABLE UserInfo (
  user_id serial,
  country varchar(128),
  display_name varchar(128),
  email varchar(128),
  explicit_content boolean,
  followers int,
  image varchar(256),
  product varchar(128),
  type varchar(128),
  uri varchar(128),
  PRIMARY KEY(user_id)
);

CREATE TABLE AuthInfo (
  user_id int REFERENCES UserInfo(user_id),
  access_token varchar(128) not null,
  token_type varchar(128) not null,
  scope varchar(128) not null,
  expires_in int not null,
  refresh_token varchar(128) not null
);

INSERT INTO UserInfo (
  user_id,
  country,
  display_name,
  email,
  explicit_content,
  followers,
  image,
  product,
  type,
  uri
  )
VALUES (
  1,
  'test_Country',
  'Test User',
  'test@test.com',
  true,
  15,
  'test_img_url',
  'test_product',
  'test_type', 
  'test_uri'
);

INSERT INTO AuthInfo 
  (user_id, access_token, token_type, scope, expires_in, refresh_token) 
VALUES 
  (1, 'thisIsATestToken', 'testTokenType', 'testTokenScope', 1337, 'testRefreshToken');
