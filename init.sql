CREATE DATABASE IF NOT EXISTS h9homework;
	
USE h9homework;
	
CREATE TABLE IF NOT EXISTS users (
	id BIGINT NOT NULL UNIQUE,
	username VARCHAR(40) NOT NULL UNIQUE,
	password VARCHAR(40) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS wallets (
	id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
	user_id BIGINT NOT NULL,
	name VARCHAR(40) NOT NULL,
	amount BIGINT,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
	
INSERT users (id, username, password)
VALUES (0, "admin", "pass");