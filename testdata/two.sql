-- qry: DeleteUsersByIds
DELETE FROM `users` WHERE `user_id` IN ({ids});
