-- Create test database
CREATE DATABASE IF NOT EXISTS test;
USE test;

-- Drop table if exists
DROP TABLE IF EXISTS t;

-- Create simple table with single integer column
CREATE TABLE t (
    i INT UNSIGNED NOT NULL,
    PRIMARY KEY (i)
) ENGINE=InnoDB;

-- Drop procedure if exists
DROP PROCEDURE IF EXISTS insert_million_shuffled;

DELIMITER $$

CREATE PROCEDURE insert_million_shuffled()
BEGIN
    -- Create temporary table with numbers 1..1,000,000
    CREATE TEMPORARY TABLE temp_numbers (
        num INT UNSIGNED NOT NULL,
        PRIMARY KEY (num)
    ) ENGINE=InnoDB;

    SELECT 'Generating numbers 1..1,000,000...' AS info;

    -- Generate 1,000,000 rows using set-based math (FAST)
    INSERT INTO temp_numbers (num)
    SELECT
        a.n
      + b.n * 10
      + c.n * 100
      + d.n * 1000
      + e.n * 10000
      + f.n * 100000
      + 1
    FROM
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) a
    CROSS JOIN
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) b
    CROSS JOIN
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) c
    CROSS JOIN
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) d
    CROSS JOIN
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) e
    CROSS JOIN
        (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
         UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) f
    WHERE
        a.n
      + b.n * 10
      + c.n * 100
      + d.n * 1000
      + e.n * 10000
      + f.n * 100000
      < 1000000;

    SELECT 'Shuffling and inserting 1,000,000 rows...' AS info;

    -- Insert in random order
    INSERT INTO t (i)
    SELECT num
    FROM temp_numbers
    ORDER BY RAND();

    DROP TEMPORARY TABLE temp_numbers;

    SELECT 'Insert complete!' AS info;
END$$

DELIMITER ;

-- Run the procedure
SELECT 'Starting to insert 1,000,000 shuffled records...' AS info;
CALL insert_million_shuffled();

-- Clean up
DROP PROCEDURE insert_million_shuffled;

-- Final stats
SELECT COUNT(*) AS total_rows FROM t;

SELECT 'Sample data (first 10 rows by insertion order):' AS info;
SELECT * FROM t LIMIT 10;

SELECT 'Sample data (first 10 rows sorted):' AS info;
SELECT * FROM t ORDER BY i LIMIT 10;
