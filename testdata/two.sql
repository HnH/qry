-- qry: DeleteUsersByIds
DELETE FROM `users` WHERE `user_id` IN ({ids});

-- qry: UglyMultiLineQuery
  SELECT * FROM
`users` WHERE
    YEAR(`birth_date`) >
2000;
