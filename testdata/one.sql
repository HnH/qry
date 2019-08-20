-- qry: InsertUser
INSERT INTO `users` (`name`) VALUES (?);

-- qry: GetUserById
SELECT * FROM `users` WHERE `user_id` = ?;
