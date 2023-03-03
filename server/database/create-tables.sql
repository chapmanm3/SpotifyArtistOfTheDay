DROP TABLE IF EXISTS authInfo;

CREATE TABLE authInfo (
  id serial primary key,
  access_token varchar(128) not null,
  token_type varchar(128) not null,
  scope varchar(128) not null,
  expires_in int not null,
  refresh_token varchar(128) not null
);

INSERT INTO authInfo 
  (access_token, token_type, scope, expires_in, refresh_token) 
VALUES 
  ('thisIsATestToken', 'testTokenType', 'testTokenScope', 1337, 'testRefreshToken');
