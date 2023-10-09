CREATE FUNCTION functest_S_1(a text, b date) RETURNS boolean
    LANGUAGE SQL
    RETURN a = 'abcd' AND b > '2001-01-01';
CREATE FUNCTION functest_S_2(a text[]) RETURNS int
    RETURN a[1]::int;
CREATE FUNCTION functest_S_3() RETURNS boolean
    RETURN false;
CREATE FUNCTION functest_S_10(a text, b date) RETURNS boolean
    LANGUAGE SQL
BEGIN ATOMIC
SELECT a = 'abcd' AND b > '2001-01-01';
END;
CREATE FUNCTION functest_S_13() RETURNS boolean
BEGIN ATOMIC
SELECT 1;
SELECT false;
END;
CREATE FUNCTION functest_S_15(x int) RETURNS boolean
    LANGUAGE SQL
BEGIN ATOMIC
select case when x % 2 = 0 then true else false end; -- tricky parsing
END;
CREATE FUNCTION functest_sri2() RETURNS SETOF int
LANGUAGE SQL
STABLE
BEGIN ATOMIC
SELECT * FROM functest3;
END;
CREATE TABLE functest1 (i int);
CREATE FUNCTION functest_S_16(a int, b int) RETURNS void
    LANGUAGE SQL
BEGIN ATOMIC
INSERT INTO functest1 SELECT a + $2;
END;

CREATE FUNCTION functest_S_15(x int) RETURNS boolean
    LANGUAGE SQL
BEGIN ATOMIC
select case when x % 2 = 0 then true else false end;
select case when x % 2 = 0 then true else false end;
select case when x % 2 = 0 then true else false end;
select case when x % 2 = 0 then true else false end;
END;

CREATE FUNCTION "functest_s_15" ("x" integer) RETURNS boolean LANGUAGE SQL BEGIN ATOMIC
SELECT
         CASE
             WHEN ((x % 2) = 0) THEN true
             ELSE false
         END AS "case";
END;