DROP TABLE IF EXISTS authInfo;
DROP TABLE IF EXISTS userInfo;

CREATE TABLE userInfo (
  user_id serial primary key,
  country varchar(128) not null,
  display_name varchar(128) not null,
  email varchar(128) not null,
  explicit_content boolean not null,
  followers int not null,
  image varchar(128) not null,
  product varchar(128) not null,
  type varchar(128) not null,
  uri varchar(128) not null,
);

CREATE TABLE authInfo (
  user_id int REFERENCES userInfo(user_id),
  access_token varchar(128) not null,
  token_type varchar(128) not null,
  scope varchar(128) not null,
  expires_in int not null,
  refresh_token varchar(128) not null
);

INSERT INTO userInfo (
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

INSERT INTO authInfo 
  (user_id, access_token, token_type, scope, expires_in, refresh_token) 
VALUES 
  (1, 'thisIsATestToken', 'testTokenType', 'testTokenScope', 1337, 'testRefreshToken');
