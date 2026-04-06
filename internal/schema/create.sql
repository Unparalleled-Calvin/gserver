-- table for use login and authentication
CREATE TABLE IF NOT EXISTS `user` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(30) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL COMMENT 'password hashed by bcrypt'
)