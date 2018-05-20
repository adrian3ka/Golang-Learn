My SQL Table Structure

CREATE DATABASE gophr;
CREATE TABLE `images` (
 `id` varchar(255) NOT NULL DEFAULT '',
 `user_id` varchar(255) NOT NULL,
 `name` varchar(255) NOT NULL DEFAULT '',
 `location` varchar(255) NOT NULL DEFAULT '',
 `description` text NOT NULL,
 `size` int(11) NOT NULL,
 `created_at` datetime NOT NULL,
 PRIMARY KEY (`id`),
 KEY `user_id_idx` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


For change user and password in Database, try to edit this line in mysql.go
db, err := NewMySQLDB("username:password@tcp(127.0.0.1:3306)/databasename")