USE GOPRAC;
-- 基本CRUD操作
-- create table students
CREATE TABLE IF NOT EXISTS `students` (
	`id` int not null auto_increment,
  `name` varchar(255),
  `age` int,
  `grade` varchar(255),
  primary key(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- insert "张三"
INSERT INTO students (`name`, `age`, `grade`) VALUE ("张三", 20, "三年级");
-- select age > 18
SELECT * FROM students WHERE age > 18;
-- update "张三"
SET SQL_SAFE_UPDATES = 0;
UPDATE students SET grade = "四年级" WHERE name = "张三";
SET SQL_SAFE_UPDATES = 1;
-- delete age < 15
SET SQL_SAFE_UPDATES = 0;
DELETE FROM students WHERE age < 15;
SET SQL_SAFE_UPDATES = 1;

-- 事务语句
CREATE TABLE IF NOT EXISTS `accounts` (
  `id` int not null auto_increment,
  primary key(`id`),
  `balance` int DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` int not null auto_increment,
  primary key(`id`),
  `from_account_id` int not null,
  `to_account_id` int not null,
  `amount` int not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT INTO accounts (balance) VALUE (300);
INSERT INTO accounts (balance) VALUE (500);
-- 1账户向2账户转账100
BEGIN;
UPDATE accounts SET balance=balance-100 WHERE id=1;
UPDATE accounts SET balance=balance+100 WHERE id=2;
INSERT INTO transactions (from_account_id, to_account_id, amount) VALUE (1, 2, 100);
COMMIT;
