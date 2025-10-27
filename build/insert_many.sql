
DO $$
DECLARE
    i INTEGER := 1;
    max_rows INTEGER := 200000;
BEGIN
    LOOP
        INSERT INTO texts(title, description, price)
        VALUES ('Component ' || i,'Description'||i,  i+100);
        i := i + 1;
        EXIT WHEN i > max_rows;
    END LOOP;
END $$;
