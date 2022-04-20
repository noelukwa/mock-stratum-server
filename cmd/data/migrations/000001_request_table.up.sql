
CREATE TABLE IF NOT EXISTS requests (
  method VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  id VARCHAR(255) NOT NULL,
  user_id VARCHAR(255), 
);


CREATE TABLE IF NOT EXISTS subscriptions (
    req_id VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL,
    id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    user_agent VARCHAR(255)
);