-- qry: EscapedJSONQuery
INSERT INTO "data" (id, "data") VALUES
  (1, '{"test": 1}'),
  (2, '{"test": 2}');

-- qry: EscapedByteaQuery
INSERT INTO bin (id, "data") VALUES
  (1, E'\\x3aaab6e7fb7245cc9785653a0d9ffc4a5ce0f974');
