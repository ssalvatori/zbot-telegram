CREATE TABLE `users` (
  `id` int(3) NOT NULL AUTO_INCREMENT,
  `username` varchar(200) NOT NULL,
  `level` int(1) NOT NULL DEFAULT '10',
  `chat` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`,`chat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `ignore_list` (
  `id` int(3) NOT NULL AUTO_INCREMENT,
  `username` varchar(200) NOT NULL,
  `chat` varchar(200) NOT NULL,
  `since` int(11) unsigned NOT NULL,
  `until` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`,`chat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `definitions` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `term` varchar(250) NOT NULL,
  `meaning` text NOT NULL,
  `chat` varchar(200) NOT NULL,
  `author` varchar(200) NOT NULL,
  `locked` tinyint(11) NOT NULL DEFAULT '0',
  `active` tinyint(11) NOT NULL DEFAULT '1',
  `hits` int(11) DEFAULT '0',
  `link` int(5) DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `locked_by` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `term_chat` (`term`,`chat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;