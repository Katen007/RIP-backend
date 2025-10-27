
EXPLAIN ANALYZE SELECT * FROM texts WHERE id BETWEEN 10 AND 10000;
SELECT * from pg_indexes WHERE tablename = 'texts';
CREATE INDEX IF NOT EXISTS texts_title_idx ON texts(title);
DROP INDEX IF EXISTS texts_title_idx;
ANALYSE;
EXPLAIN ANALYZE SELECT * FROM texts WHERE title IN ('Text 100000', 'Text 20000');
SELECT COUNT(*) FROM texts;
