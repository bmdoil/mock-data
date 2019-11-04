SELECT * FROM (  
		       SELECT n.nspname || '.' || c.relname tablename,  
		       	con.conname conname,  
		              pg_catalog.pg_get_constraintdef(con.oid, true) concol, 
		              'constraint' contype  
		       FROM  pg_catalog.pg_class c,  
		       	  pg_catalog.pg_constraint con,  
		       	  pg_namespace n  
		       WHERE conrelid = c.oid  
		       AND n.oid = c.relnamespace  
		       AND contype IN ('u','f','c','p')  
		       UNION  
		       SELECT schemaname || '.' || tablename tablename,  
		       	   indexname conname,  
		                 indexdef concol,  
		       	   'index' contype  
		       FROM   pg_indexes   
		       WHERE  schemaname IN (SELECT nspname   
		       FROM   pg_namespace   
		       WHERE  nspname NOT IN (   
		       'pg_catalog',   
		       'information_schema',  
		       'pg_aoseg',  
		       'gp_toolkit',  
		       'pg_toast', 'pg_bitmapindex' ))   
		       AND indexdef LIKE 'CREATE UNIQUE%'   
		) a ORDER BY contype




SELECT n.nspname || '.' || c.relname tablename,  
				   con.conname constraint_name, 
			       pg_catalog.pg_get_constraintdef(con.oid, true) constraint 
			FROM  pg_catalog.pg_class c, 
				  pg_catalog.pg_constraint con, 
				  pg_namespace n 
			WHERE conrelid = c.oid 
			AND n.oid = c.relnamespace 
			AND contype = 'c' 
			ORDER  BY tablename ;


-- Get all tables in database

SELECT n.nspname    || '.' || c.relname   
		FROM   pg_catalog.pg_class c  
		       LEFT JOIN pg_catalog.pg_namespace n  
		              ON n.oid = c.relnamespace  
		WHERE  c.relkind IN ( 'r', '' )  
		       AND n.nspname <> 'pg_catalog'  
		       AND n.nspname <> 'information_schema'  
		       AND n.nspname !~ '^pg_toast'  
		       AND n.nspname <> 'gp_toolkit'  
		       AND c.relkind = 'r'  
		       AND c.relstorage IN ('a', 'h')  
		ORDER  BY 1;


 SELECT   a.attname,  
		        pg_catalog.Format_type(a.atttypid, a.atttypmod),  
			 COALESCE((SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128)  
		        FROM pg_catalog.pg_attrdef d  
		        WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef), '')  
		FROM     pg_catalog.pg_attribute a  
		WHERE    a.attrelid = 'sales'::regclass  
		AND      a.attnum > 0  
		AND      NOT a.attisdropped  
		ORDER BY a.attnum ;





SELECT c.oid,
  n.nspname,
  c.relname
FROM pg_catalog.pg_class c
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relname ~ '^(sales)$'
  AND pg_catalog.pg_table_is_visible(c.oid)
ORDER BY 2, 3;

SELECT relchecks, relkind, relhasindex, relhasrules, reltriggers <> 0, relhasoids, '', reltablespace, relstorage
FROM pg_catalog.pg_class WHERE oid = '16385';

SELECT a.attname,
  pg_catalog.format_type(a.atttypid, a.atttypmod),
  (SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128)
   FROM pg_catalog.pg_attrdef d
   WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef),
  a.attnotnull, a.attnum
FROM pg_catalog.pg_attribute a;



SELECT partitionschemaname schema, 
		partitiontablename tablename, 
		partitiontype type, 
		partitionrangestart rangestart, 
		partitionstartinclusive startinclusive, 
		partitionrangeend rangeend, 
		partitionendinclusive endinclusive,
		partitioneveryclause everyclause 
		FROM	 pg_catalog.pg_partitions

SELECT partitionschemaname || '.' || partitiontablename
		FROM	 pg_catalog.pg_partitions

		
CREATE TABLE sales_2 (trans_id int, date date, amount 
decimal(9,2), region text) 
DISTRIBUTED BY (trans_id)
PARTITION BY RANGE (date)
SUBPARTITION BY LIST (region)
SUBPARTITION TEMPLATE
( SUBPARTITION usa VALUES ('usa'), 
  SUBPARTITION asia VALUES ('asia'), 
  SUBPARTITION europe VALUES ('europe'), 
  DEFAULT SUBPARTITION other_regions)
  (START (date '2011-01-01') INCLUSIVE
   END (date '2012-01-01') EXCLUSIVE
   EVERY (INTERVAL '1 month'), 
   DEFAULT PARTITION outlying_dates );


//Get all range check constraints


   SELECT p.partitionschemaname schema,
   p.partitiontablename relname,
   c.conname, 
   a.attname colname,  
   TRIM (TRAILING '::date' from CAST (p.partitionrangestart as TEXT)) rangestart, 
   TRIM (TRAILING '::date' from CAST (p.partitionrangeend as TEXT)) rangeend, 
   p.partitionstartinclusive startinclusive, 
   p.partitionendinclusive endinclusive
    FROM pg_class cl 
    JOIN pg_partitions p on cl.relname=p.partitiontablename
    JOIN pg_attribute a on cl.oid=a.attrelid
	JOIN pg_constraint c on cl.oid=c.conrelid
   WHERE a.attnum = ANY (c.conkey) and c.contype='c' and p.partitiontype='range'
   ORDER BY conname desc;


//Get all list check constraints

 SELECT p.partitionschemaname schema,
   p.partitiontablename relname, 
   c.conname, 
   a.attname colname
    FROM pg_class cl 
    JOIN pg_partitions p on cl.relname=p.partitiontablename
    JOIN pg_attribute a on cl.oid=a.attrelid
	JOIN pg_constraint c on cl.oid=c.conrelid
   WHERE a.attnum = ANY (c.conkey) and c.contype='c' and p.partitiontype='list'
   ORDER BY conname desc;


-- Get all check constraints
   SELECT p.partitionschemaname || '.' || p.partitiontablename relname,
   c.conname, 
   a.attname colname,  
   TRIM (TRAILING '::date' from CAST (p.partitionrangestart as TEXT)) rangestart, 
   TRIM (TRAILING '::date' from CAST (p.partitionrangeend as TEXT)) rangeend, 
   p.partitionstartinclusive startinclusive, 
   p.partitionendinclusive endinclusive
    FROM pg_class cl 
    JOIN pg_partitions p on cl.relname=p.partitiontablename
    JOIN pg_attribute a on cl.oid=a.attrelid
	JOIN pg_constraint c on cl.oid=c.conrelid
   WHERE a.attnum = ANY (c.conkey) and c.contype='c'
   ORDER BY conname desc;



-- Get all tables from database, and if table is a partition table
SELECT all_tables.relname, case when part_tables.partition is null then 'no' else part_tables.partition end as partitiontable
FROM
(
    SELECT n.nspname    || '.' || c.relname   relname 
	FROM   pg_catalog.pg_class c  
	       LEFT JOIN pg_catalog.pg_namespace n  
	              ON n.oid = c.relnamespace  
	WHERE  c.relkind IN ( 'r', '' )  
	       AND n.nspname <> 'pg_catalog'  
	       AND n.nspname <> 'information_schema'  
	       AND n.nspname !~ '^pg_toast'  
	       AND n.nspname <> 'gp_toolkit'  
	       AND c.relkind = 'r'  
	       AND c.relstorage IN ('a', 'h')  
	ORDER  BY 1
) all_tables
LEFT JOIN (
SELECT partitionschemaname || '.' || partitiontablename as relname, 'child' as partition    
FROM pg_catalog.pg_partitions 
UNION SELECT schemaname || '.' || tablename as tbl, 'parent' as partition 
FROM pg_catalog.pg_partitions
) part_tables
ON all_tables.relname=part_tables.relname;
 

// All partition tables
WITH partition_tables as (
	SELECT partitionschemaname || '.' || partitiontablename as tbl, 'child' as partition    
	FROM pg_catalog.pg_partitions 
	UNION SELECT DISTINCT schemaname || '.' || tablename as tbl, 'parent' as partition 
	FROM pg_catalog.pg_partitions), 
 all_tables as (   SELECT n.nspname    || '.' || c.relname as relname
		FROM   pg_catalog.pg_class c  
		       LEFT JOIN pg_catalog.pg_namespace n  
		              ON n.oid = c.relnamespace  
		WHERE  c.relkind IN ( 'r', '' )  
		       AND n.nspname <> 'pg_catalog'  
		       AND n.nspname <> 'information_schema'  
		       AND n.nspname !~ '^pg_toast'  
		       AND n.nspname <> 'gp_toolkit'  
		       AND c.relkind = 'r'  
		       AND c.relstorage IN ('a', 'h')  
		ORDER  BY 1)
SELECT * FROM partition_tables 
LEFT JOIN all_tables on all_tables.relname=partition_tables.tbl;

partition_tables where partitiontype='range';

SELECT partitionschemaname || '.' || partitiontablename as tbl, 'child' as partition    
FROM pg_catalog.pg_partitions UNION SELECT DISTINCT schemaname || '.' || tablename as tbl, 'parent' as partition from pg_catalog.pg_partitions
UNION SELECT DISTINCT 
                                                                                                                           


 case when all_tables.relname=part_tables.tbl then 'true' else 'false' end as
   FROM 

   SELECT n.nspname    || '.' || c.relname   
		FROM   pg_catalog.pg_class c  
		       LEFT JOIN pg_catalog.pg_namespace n  
		              ON n.oid = c.relnamespace  
		WHERE  c.relkind IN ( 'r', '' )  
		       AND n.nspname <> 'pg_catalog'  
		       AND n.nspname <> 'information_schema'  
		       AND n.nspname !~ '^pg_toast'  
		       AND n.nspname <> 'gp_toolkit'  
		       AND c.relkind = 'r'  
		       AND c.relstorage IN ('a', 'h')  
		ORDER  BY 1;