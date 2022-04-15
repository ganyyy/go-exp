# 作用

1. 方便key的导入/导出
2. 线上出现过因为备库的redis/ardb版本和主库的版本不一致的情况, 导致无法实现数据的restore/dump, 这是一种通用的数据导入/导出工具

# 使用方法

### 将线上数据导入到本地

```shell

  # 服务器

  ./redis-key-backup dump --host=localhost:6379 --auth="" --db=3 --key=hash2 --file=output.json
 
  # 本地
  ./redis-key-backup restore --host=localhost:6379 --db=4 --key=hash4 --file=output.json
```

### 线上数据备库 -> 主库

```shell
  # 适用于脚本命令行的批量导入
  ./redis-key-backup dump --host=localhost:6379 --db=3 --key=hash2 --output | ./redis-key-backup restore --host=localhost:6379 --db=4 --key=hash4 --input

```