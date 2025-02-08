# usage
```
go run main.go --host 192.168.219.1 --pwd 11111
go run dmTable.go -host 192.168.1.1 -user test -pwd 111111111 -total 100000 -parallel 10
```

# sql
```
SELECT '/*aa*/ddrop table '|| table_name||' ;' FROM user_tables;

/*aa*/drop table new_thread684_table_86;

CREATE tablespace test DATAFILE 'test.dbf' SIZE 200 autoextend ON NEXT 200;
CREATE USER test IDENTIFIED BY 111111111 DEFAULT tablespace test;
GRANT dba TO test;

SELECT * FROM dba_tables WHERE owner='test';

```
