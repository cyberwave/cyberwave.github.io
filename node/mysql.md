# 一、修改表结构

示例表结构：
```sql
CREATE TABLE IF NOT EXISTS 
        student (id BIGINT(20) NOT NULL AUTO_INCREMENT,
        name VARCHAR(64) NOT NULL UNIQUE,
        PRIMARY KEY(id)) 
        DEFAULT CHARSET=utf8  COLLATE utf8_bin;`
```

1. 查看表结构
```sql
desc student
```

2. 查看建表语句
```sql
show create table student
```

3. 添加一列
```sql
alter table student add column age int not null -- 默认在最后一列添加
alter table student add column collage varchcar(54) after name; -- 在某一列后面添加一列
alter table student add column grade varchar(64) not null first; -- 添加一列到第一列
```

4. 删除列
```sql
alter table student drop column age;
```

5. 修改列
把 age 改成 class
```sql
alter table student change age class varchar(10);
```

6. 删除主键
```sql
alter table student drop primary key;
```

7. 修改表名
```sql
rename table student to students
```

8. 复制一列
```sql
update student set age = grade
```


