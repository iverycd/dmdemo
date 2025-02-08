package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/godoes/gorm-dameng"
	"gorm.io/gorm"
	"math"
	"strconv"
	"sync"
	"time"
)

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// todo 待增加命令行参数 并行数，总共要创建的表总数，总共要创建的表总数/并行数=每个协程for循环的次数
var User, Host, Pwd, Total, Parallel arrayFlags
var GlobalCount = 0
var mutex sync.Mutex
var wg sync.WaitGroup // 创建WaitGroup实例

func init() {
	// 从命令行获取flag参数
	flag.Var(&User, "user", "db_user")
	flag.Var(&Host, "host", "ip")
	flag.Var(&Pwd, "pwd", "password")
	flag.Var(&Total, "total", "total table")
	flag.Var(&Parallel, "parallel", "parallel thread")
	// 通过Parse解析获取到配置文件
	flag.Parse()
	if len(Host) == 0 || len(Pwd) == 0 || len(User) == 0 {
		panic(errors.New("config file not specified."))
	}

}

//func createTable(db *gorm.DB, index int) {
//	count := 1
//	for {
//		defer wg.Done() // 通知WaitGroup当前goroutine已结束
//		mutex.Lock()
//		count++
//		GlobalCount++
//		fmt.Println("GlobalCount:", GlobalCount)
//		numTotal, _ := strconv.Atoi(Total[0])
//		sql := fmt.Sprintf("create table if not exists new_thread%d_table_%d(id1 int,id2 int,id3 int,id4 int,id5 int,name1 varchar(20),name2 varchar(20),name3 varchar(20),name4 varchar(20),name5 varchar(20))", index, GlobalCount)
//		db.Exec(sql)
//		if GlobalCount == numTotal {
//			fmt.Println("当前已创建表:", numTotal, "程序退出")
//			os.Exit(0)
//		}
//		mutex.Unlock()
//
//		//fmt.Println("线程", index, "第", count, "次执行建表-", sql)
//	}
//}

func createTable(db *gorm.DB, index int, count int) {
	for i := 1; i < count+1; i++ {
		fmt.Println("当前协程-", index, " 创建第", i, "张表")
		// 10列
		//sql := fmt.Sprintf("create table if not exists new_thread%d_table_%d(id1 int,id2 int,id3 int,id4 int,id5 int,name1 varchar(20),name2 varchar(20),name3 varchar(20),name4 varchar(20),name5 varchar(20))", index, i)
		// 20列
		sql := fmt.Sprintf("create table if not exists new_thread%d_table_%d(id1 int,id2 int,id3 int,id4 int,id5 int,id6 int,id7 int,id8 int,id9 int,id10 int,name1 varchar(20),name2 varchar(20),name3 varchar(20),name4 varchar(20),name5 varchar(20),name6 varchar(20),name7 varchar(20),name8 varchar(20),name9 varchar(20),name10 varchar(20))", index, i)
		db.Exec(sql)
	}
	defer wg.Done() // 通知WaitGroup当前goroutine已结束
}

func main() {
	startTime := time.Now()
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
	// 并行度
	parallelNum, _ := strconv.Atoi(Parallel[0])
	// 目标创建总数
	totalNum, _ := strconv.Atoi(Total[0])
	// 每个协程循环次数即每个协程创建的表总数
	countNum := math.Ceil(float64(totalNum / parallelNum))
	// 设置等待的goroutine数量
	wg.Add(parallelNum)

	for i := 0; i < parallelNum; i++ {
		//fmt.Println("并发执行", i, time.Now())
		//fmt.Println("GlobalCount:", GlobalCount)
		//if GlobalCount == 10 {
		//	os.Exit(0)
		//}
		// db连接,协程号,每个协程循环的次数
		go createTable(db, i, int(countNum))
	}

	//time.Sleep(5 * time.Second)
	wg.Wait() // 等待所有goroutine完成，主goroutine将继续执行到这
	endTime := time.Now()
	fmt.Println("用时:", endTime.Sub(startTime).Seconds(), "秒")

}

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
