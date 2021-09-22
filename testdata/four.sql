
-- qry: QueryWithComments
SELECT name, email
FROM users
INNER JOIN user_details -- comment
ON  users.id = user_details.user_id /* another comment */
/*
  multiline
  comment
*/
WHERE users.id = $1;