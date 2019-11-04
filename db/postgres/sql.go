package postgres

import "strings"

// Postgres version
func PGVersion() string {
	return "select version()"
}

// Get all table query
// Postgres 9 and above
func PGAllTablesQry1() string {
	return "SELECT n.nspname    || '.' || c.relname  " +
		"FROM   pg_catalog.pg_class c " +
		"       LEFT JOIN pg_catalog.pg_namespace n " +
		"              ON n.oid = c.relnamespace " +
		"WHERE  c.relkind IN ( 'r', '' ) " +
		"       AND n.nspname <> 'pg_catalog' " +
		"       AND n.nspname <> 'information_schema' " +
		"       AND n.nspname !~ '^pg_toast' " +
		"       AND n.nspname !~ '^gp_toolkit' " +
		"       AND c.relkind = 'r' " +
		"ORDER  BY 1 "
}

// Greenplum, HDB , postgres 8.3 etc
func PGAllTablesQry2() string {
	return "SELECT n.nspname    || '.' || c.relname  " +
		"FROM   pg_catalog.pg_class c " +
		"       LEFT JOIN pg_catalog.pg_namespace n " +
		"              ON n.oid = c.relnamespace " +
		"WHERE  c.relkind IN ( 'r', '' ) " +
		"       AND n.nspname <> 'pg_catalog' " +
		"       AND n.nspname <> 'information_schema' " +
		"       AND n.nspname !~ '^pg_toast' " +
		"       AND n.nspname <> 'gp_toolkit' " +
		"       AND c.relkind = 'r' " +
		"       AND c.relstorage IN ('a', 'h') " +
		"ORDER  BY 1 "
}

// Get all tables from database, and if table is a partition table
func GPAllTablesQryPartitions() string {
	return " SELECT all_tables.relname, case when part_tables.partition is null then 'no' else part_tables.partition end as partitiontable " +
		" FROM ( " +
		" SELECT n.nspname    || '.' || c.relname   relname " +
		" FROM   pg_catalog.pg_class c " +
		" LEFT JOIN pg_catalog.pg_namespace n " +
		" ON n.oid = c.relnamespace " +
		" WHERE  c.relkind IN ( 'r', '' ) " +
		"	AND n.nspname <> 'pg_catalog' " +
		"	AND n.nspname <> 'information_schema' " +
		"	AND n.nspname !~ '^pg_toast' " +
		"	AND n.nspname <> 'gp_toolkit' " +
		"	AND c.relkind = 'r' " +
		"	AND c.relstorage IN ('a', 'h') " +
		" ORDER  BY 1) all_tables " +
		" LEFT JOIN ( " +
		" 	SELECT partitionschemaname || '.' || partitiontablename as relname, 'child' as partition " +
		" 	FROM pg_catalog.pg_partitions " +
		" 	UNION SELECT schemaname || '.' || tablename as tbl, 'parent' as partition " +
		" 	FROM pg_catalog.pg_partitions) part_tables " +
		" ON all_tables.relname=part_tables.relname "
}

func GPTableQryPartitions(table string) string {
	return " SELECT all_tables.relname, case when part_tables.partition is null then 'no' else part_tables.partition end as partitiontable " +
		" FROM ( " +
		" SELECT n.nspname    || '.' || c.relname   relname " +
		" FROM   pg_catalog.pg_class c " +
		" LEFT JOIN pg_catalog.pg_namespace n " +
		" ON n.oid = c.relnamespace " +
		" WHERE  c.relkind IN ( 'r', '' ) " +
		"	AND n.nspname <> 'pg_catalog' " +
		"	AND n.nspname <> 'information_schema' " +
		"	AND n.nspname !~ '^pg_toast' " +
		"	AND n.nspname <> 'gp_toolkit' " +
		"	AND c.relkind = 'r' " +
		"	AND c.relstorage IN ('a', 'h') " +
		"   AND c.oid='" + table + "'::regclass" +
		" ORDER  BY 1) all_tables " +
		" LEFT JOIN ( " +
		" 	SELECT partitionschemaname || '.' || partitiontablename as relname, 'child' as partition " +
		" 	FROM pg_catalog.pg_partitions " +
		" 	UNION SELECT schemaname || '.' || tablename as tbl, 'parent' as partition " +
		" 	FROM pg_catalog.pg_partitions) part_tables " +
		" ON all_tables.relname=part_tables.relname"
}

// Get all columns from the table query
// Postgres 9 and above
func PGColumnQry1(table string) string {
	return "SELECT   a.attname, " +
		"        pg_catalog.Format_type(a.atttypid, a.atttypmod), " +
		"	 COALESCE((SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128) " +
		"        FROM pg_catalog.pg_attrdef d " +
		"        WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef), '') " +
		"FROM     pg_catalog.pg_attribute a " +
		"WHERE    a.attrelid = '" + table + "'::regclass " +
		"AND      a.attnum > 0 " +
		"AND      NOT a.attisdropped " +
		"ORDER BY a.attnum "
}

// Postgres 8.3, GPDB, HDB etc
func PGColumnQry2(table string) string {
	return "SELECT         a.attname, " +
		"               pg_catalog.Format_type(a.atttypid, a.atttypmod), " +
		"	        COALESCE((SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128) " +
		"                FROM pg_catalog.pg_attrdef d " +
		"                WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef), '') " +
		"FROM            pg_catalog.pg_attribute a " +
		"LEFT OUTER JOIN pg_catalog.pg_attribute_encoding e " +
		"ON              e.attrelid = a .attrelid " +
		"AND             e.attnum = a.attnum " +
		"WHERE           a.attrelid = '" + table + "'::regclass " +
		"AND             a.attnum > 0 " +
		"AND             NOT a.attisdropped " +
		"ORDER BY        a.attnum"
}

// Get list of partition tables in DB
func GetPartitionTables() string {
	return "SELECT partitionschemaname schema," +
		"partitiontablename tablename," +
		"partitiontype type," +
		"partitionrangestart rangestart," +
		"partitionstartinclusive startinclusive," +
		"partitionrangeend rangeend," +
		"partitionendinclusive endinclusive," +
		"partitioneveryclause everyclause," +
		"FROM	 pg_catalog.pg_partitions"
}
func GetAllCheckConstraints() string {
	return "SELECT p.partitionschemaname || '.' || p.partitiontablename relname, " +
		"c.conname, " +
		"p.partitiontype, " +
		"a.attname colname, " +
		"TRIM (TRAILING '::date' from CAST (p.partitionrangestart as TEXT)) rangestart, " +
		"TRIM (TRAILING '::date' from CAST (p.partitionrangeend as TEXT)) rangeend, " +
		"p.partitionstartinclusive startinclusive, " +
		"p.partitionendinclusive endinclusive " +
		"FROM pg_class cl " +
		"JOIN pg_partitions p on cl.relname=p.partitiontablename " +
		"JOIN pg_attribute a on cl.oid=a.attrelid " +
		"JOIN pg_constraint c on cl.oid=c.conrelid " +
		"WHERE a.attnum = ANY (c.conkey) and c.contype='c' " +
		"ORDER BY conname desc;"
}

/*
schema
relname
conname
colname
rangestart
rangeend
startinclusive
endinclusive
*/
func GetRangeConstraints() string {
	return "   SELECT p.partitionschemaname schema, " +
		"	p.partitiontablename relname, " +
		"	c.conname, " +
		"	a.attname colname, " +
		"	TRIM (TRAILING '::date' from CAST (p.partitionrangestart as TEXT)) rangestart, " +
		"	TRIM (TRAILING '::date' from CAST (p.partitionrangeend as TEXT)) rangeend, " +
		"	p.partitionstartinclusive startinclusive, " +
		"	p.partitionendinclusive endinclusive " +
		"	FROM pg_class cl " +
		"   JOIN pg_partitions p on cl.relname=p.partitiontablename " +
		"   JOIN pg_attribute a on cl.oid=a.attrelid " +
		"   JOIN pg_constraint c on cl.oid=c.conrelid " +
		"   WHERE a.attnum = ANY (c.conkey) and c.contype='c' and p.partitiontype='range'"
}

/*
schema
relname
conname
colname
*/
func GetListConstraints() string {
	return "   SELECT p.partitionschemaname schema, " +
		"	p.partitiontablename relname, " +
		"	c.conname, " +
		"	a.attname colname, " +
		"	FROM pg_class cl " +
		"   JOIN pg_partitions p on cl.relname=p.partitiontablename " +
		"   JOIN pg_attribute a on cl.oid=a.attrelid " +
		"   JOIN pg_constraint c on cl.oid=c.conrelid " +
		"   WHERE a.attnum = ANY (c.conkey)" +
		"   AND c.contype='c'" +
		"   AND p.partitiontype='list'"
}

// Save all the DDL of the constraint ( like PK(p), FK(f), CK(c), UK(u) )
func GetPGConstraintDDL(conntype string) string {
	return "	SELECT n.nspname || '.' || c.relname tablename, " +
		"		   con.conname constraint_name," +
		"	       pg_catalog.pg_get_constraintdef(con.oid, true) constriant_col" +
		"	FROM  pg_catalog.pg_class c," +
		"		  pg_catalog.pg_constraint con," +
		"		  pg_namespace n" +
		"	WHERE conrelid = c.oid" +
		"	AND n.oid = c.relnamespace" +
		"	AND contype = '" + conntype + "'" +
		"	ORDER  BY tablename "
}

// Get all the Unique index from the database
func GetPGIndexDDL() string {
	return "SELECT schemaname ||'.'|| tablename, " +
		"indexdef " +
		"FROM   pg_indexes " +
		"WHERE  schemaname IN (SELECT nspname " +
		"FROM   pg_namespace " +
		"WHERE  nspname NOT IN ( " +
		"'pg_catalog', " +
		"'information_schema'," +
		"'pg_aoseg'," +
		"'gp_toolkit'," +
		"'pg_toast', 'pg_bitmapindex' )) " +
		"AND indexdef LIKE 'CREATE UNIQUE%'"
}

// Drop statement for the table
func GetConstraintsPertab(tabname string) string {
	return "SELECT * FROM ( " +
		"       SELECT n.nspname || '.' || c.relname tablename, " +
		"       	con.conname conname, " +
		"              pg_catalog.pg_get_constraintdef(con.oid, true) concol," +
		"              'constriant' contype " +
		"       FROM  pg_catalog.pg_class c, " +
		"       	  pg_catalog.pg_constraint con, " +
		"       	  pg_namespace n " +
		"       WHERE  c.oid = '" + tabname + "'::regclass " +
		"       AND conrelid = c.oid " +
		"       AND n.oid = c.relnamespace " +
		"       AND contype IN ('u','f','c','p') " +
		"       UNION " +
		"       SELECT schemaname || '.' || tablename tablename, " +
		"       	   indexname conname, " +
		"                 indexdef concol, " +
		"       	   'index' contype " +
		"       FROM   pg_indexes  " +
		"       WHERE  schemaname IN (SELECT nspname  " +
		"       FROM   pg_namespace  " +
		"       WHERE  nspname NOT IN (  " +
		"       'pg_catalog',  " +
		"       'information_schema', " +
		"       'pg_aoseg', " +
		"       'gp_toolkit', " +
		"       'pg_toast', 'pg_bitmapindex' ))  " +
		"       AND indexdef LIKE 'CREATE UNIQUE%' " +
		"       AND schemaname || '.' || tablename = '" + tabname + "' " +
		") a ORDER BY contype" // Ensuring the constraint remains on top

}

// Get the datatype of the column
func GetDatatype(tab string, columns []string) string {
	whereClause := strings.Join(columns, "' or attname = '")
	whereClause = strings.Replace(whereClause, "attname = ' ", "attname = '", -1)
	query := "SELECT attname, pg_catalog.Format_type(atttypid, atttypmod) " +
		"FROM pg_attribute WHERE (attname = " +
		"'" + whereClause + "') AND attrelid = '" + tab + "'::regclass"
	return query
}

// Primary key violation check
func GetTotalPKViolator(tab, cols string) string {
	return "SELECT COUNT(*) FROM ( " +
		GetPKViolator(tab, cols) +
		") a "
}

// Total Primary Key violator
func GetPKViolator(tab, cols string) string {
	return " SELECT " + cols +
		" FROM " + tab +
		" GROUP BY " + cols +
		" HAVING COUNT(*) > 1 "
}

// Fix int PK violators.
func UpdateIntPKey(tab, col, dt string) string {
	return " UPDATE " + tab +
		" SET " + col + " = " + col + "+" + "trunc(random() * 100 + 1)::" + dt +
		" WHERE " + col + " IN ( " + GetPKViolator(tab, col) + " )"
}

// Fix String PK Violators
func UpdatePKey(tab, col, whichrow, newdata string) string {
	return " UPDATE " + tab +
		" SET " + col + " = '" + newdata + "'" +
		" WHERE ctid = ( SELECT ctid FROM " + tab +
		" WHERE " + col + " = '" + whichrow + "' LIMIT 1 )"
}

// Get the foriegn violations keys
func GetFKViolators(tab, col, reftab, refcol string) string {
	return "SELECT " + col + " FROM " + tab + " where " + col + " NOT IN ( SELECT " + refcol + " from " + reftab + " )"
}

// Get total FK violators
func GetTotalFKViolators(tab, col, reftab, refcol string) string {
	return "SELECT COUNT(*) FROM (" + GetFKViolators(tab, col, reftab, refcol) + ") a"
}

// Total rows of the table
func TotalRows(tab string) string {
	return "SELECT COUNT(*) FROM " + tab
}

// Update FK violators
func UpdateFKeys(fktab, fkcol, reftab, refcol, whichrow, totalRows string) string {
	return "UPDATE " + fktab + " SET " + fkcol +
		"=(SELECT " + refcol + " FROM " + reftab +
		" OFFSET floor(random()*" + totalRows + ") LIMIT 1)" +
		" WHERE " + fkcol + "='" + whichrow + "'"
}

// Check key violation check
func GetTotalCKViolator(tab, column, ckconstraint string) string {
	return "SELECT COUNT(*) FROM ( " +
		GetCKViolator(tab, column, ckconstraint) +
		") a "
}

// Check Constraint violation
func GetCKViolator(tab, column, ckconstraint string) string {
	return "SELECT " + column +
		"FROM " + tab + " WHERE not (" + ckconstraint + ")"
}
