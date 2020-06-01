CREATE TABLE `users` (
	`id` INT(3) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(200) NOT NULL,
	`level` INT(1) NOT NULL DEFAULT 10,
	PRIMARY KEY (`id`)
);

CREATE TABLE `ignore_list` (
    `id` INT(3) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(200) NOT NULL,
  `since` int(11) unsigned NOT NULL,
  `until` int(11) unsigned NOT NULL,
    PRIMARY KEY(`id`)
);

CREATE TABLE `definitions` (
  `id` int(5) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `term` varchar(250) UNIQUE NOT NULL,
  `meaning` text NOT NULL,
  `author` varchar(200) NOT NULL,
  `locked` tinyint(11) NOT NULL DEFAULT '0',
  `active` tinyint(11) NOT NULL DEFAULT '1',
  `hits` int(11) DEFAULT '0',
  `link` int(5) DEFAULT NULL,
  `created` timestamp default now(), 
  `modified` timestamp default now() on update now(),
  `locked_by` varchar(250) DEFAULT ''
);
