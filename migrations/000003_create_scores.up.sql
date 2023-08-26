CREATE TABLE Scores (
  user_id VARCHAR(128) PRIMARY KEY,
  score INT default 0,
  CONSTRAINT fk_t0_userid FOREIGN KEY fk_userid(user_id) REFERENCES users (id) ON DELETE cascade
);