package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/godoes/gorm-dameng"
	"gorm.io/gorm"
	"os"
	"strconv"
	"sync"
)

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var User, Host, Pwd, Total arrayFlags
var GlobalCount = 0
var mutex sync.Mutex
var wg sync.WaitGroup // 创建WaitGroup实例

func init() {
	// 从命令行获取flag参数
	flag.Var(&User, "user", "db_user")
	flag.Var(&Host, "host", "ip")
	flag.Var(&Pwd, "pwd", "password")
	flag.Var(&Total, "total", "total table")
	// 通过Parse解析获取到配置文件
	flag.Parse()
	if len(Host) == 0 || len(Pwd) == 0 || len(User) == 0 {
		panic(errors.New("config file not specified."))
	}

}

func createTable(db *gorm.DB, index int) {
	count := 1
	for {
		sql := fmt.Sprintf("create table if not exists new_thread%d_table_%d(id1 int,id2 int,id3 int,id4 int,id5 int,name1 varchar(20),name2 varchar(20),name3 varchar(20),name4 varchar(20),name5 varchar(20))", index, count)
		db.Exec(sql)
		defer wg.Done() // 通知WaitGroup当前goroutine已结束
		mutex.Lock()
		count++
		GlobalCount++
		fmt.Println("GlobalCount:", GlobalCount)
		numTotal, _ := strconv.Atoi(Total[0])
		if GlobalCount == numTotal {
			fmt.Println("当前已创建表:", numTotal, "程序退出")
			os.Exit(0)
		}
		mutex.Unlock()
		//fmt.Println("线程", index, "第", count, "次执行建表-", sql)

	}
}

func main() {
	options := map[string]string{
		"schema":         User[0],
		"appName":        "GORM连接达梦数据库示例",
		"connectTimeout": "30",
	}
	// dsn := fmt.Sprintf("dm://SYSDBA:%s@%s:5236?appName=MacPro&connectTimeout=3000&clientEncoding=GB18030", Pwd[0], Host[0])
	dsn := dameng.BuildUrl(User[0], Pwd[0], Host[0], 5236, options)
	db, err := gorm.Open(dameng.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread1", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread2", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread3", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread4", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread5", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(1 * time.Millisecond)
	//	count := 1
	//	for range ticker.C {
	//		sql := fmt.Sprintf("create table if not exists %s_table_%d(id int)", "thread6", count)
	//		db.Exec(sql)
	//		count++
	//		fmt.Println("第", count, "次执行:")
	//	}
	//}()

	wg.Add(100) // 设置等待的goroutine数量为10

	for i := 0; i < 100; i++ {
		//fmt.Println("并发执行", i, time.Now())
		//fmt.Println("GlobalCount:", GlobalCount)
		//if GlobalCount == 10 {
		//	os.Exit(0)
		//}
		go createTable(db, i)
	}

	//time.Sleep(5 * time.Second)
	wg.Wait() // 等待所有goroutine完成，主goroutine将继续执行到这

	//fmt.Println("GlobalCount:", GlobalCount)

}

//-- 50w张表 10个字段和20个字段却别
//
//SELECT '/*aa*/ddrop table '|| table_name||' ;' FROM user_tables;
//
///*aa*/drop table new_thread684_table_86;
//
//CREATE tablespace test DATAFILE 'test.dbf' SIZE 200 autoextend ON NEXT 200;
//CREATE USER test IDENTIFIED BY 111111111 DEFAULT tablespace test;
//GRANT dba TO test;
//
//
//SELECT * FROM dba_tables WHERE owner='test';
