CREATE TABLE authors (
  id INTEGER AUTO_INCREMENT,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  username VARCHAR(255),
  created_at TIMESTAMP,
  created_date TIMESTAMP,
  publish_date TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TABLE posts (
  id INTEGER AUTO_INCREMENT,
  subject VARCHAR(255),
  body TEXT,
  author_id INTEGER,
  created_at TIMESTAMP default CURRENT_TIMESTAMP,
  created_date TIMESTAMP,
  deleted_at TIMESTAMP,
  published_at TIMESTAMP,
  publish_date TIMESTAMP,
  updated_at TIMESTAMP,
  is_published BOOLEAN default 0,
  PRIMARY KEY (id)
);

CREATE TABLE comments (
  id INTEGER AUTO_INCREMENT,
  subject VARCHAR(255),
  body TEXT,
  author_id INTEGER,
  post_id INTEGER,
  created_at TIMESTAMP default CURRENT_TIMESTAMP,
  created_date TIMESTAMP,
  deleted_at TIMESTAMP,
  published_at TIMESTAMP,
  publish_date TIMESTAMP,
  updated_at TIMESTAMP,
  is_published BOOLEAN default 0,
  PRIMARY KEY (id)
);
