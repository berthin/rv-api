CREATE SCHEMA redventures;

CREATE TABLE `redventures`.`user` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(60) NOT NULL,
  `gravatar` VARCHAR(200) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `redventures`.`widget` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(60) NOT NULL,
  `color` VARCHAR(30) NULL DEFAULT 'NOT DESCRIBED',
  `price` DECIMAL(10,2) NULL DEFAULT 0.00,
  `melts` TINYINT(1) NULL DEFAULT 0,
  `inventory` INT NULL DEFAULT 0,
  PRIMARY KEY (`id`));

INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Widad Niklas', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8RpEBBTnUZhKxq9gHAV_8jVSKGF9E8p6cUzaWaQl8BAII_Elt');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Omran Yusuf', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8RpEBBTnUZhKxq9gHAV_8jVSKGF9E8p6cUzaWaQl8BAII_Elt');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Denis Leonor', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8RpEBBTnUZhKxq9gHAV_8jVSKGF9E8p6cUzaWaQl8BAII_Elt');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Maximiliano Olivia', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8RpEBBTnUZhKxq9gHAV_8jVSKGF9E8p6cUzaWaQl8BAII_Elt');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Fábio Natalia', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Juliane Tyr Sonia', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Aphrodite Janine Kristine', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Adrastea Leonard Ligeia', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Iphigeneia Franziska Balbino', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');
INSERT INTO `redventures`.`user` (`name`, `gravatar`) VALUES ('Adrastea Reinhold Cecilia', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q');

INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `inventory`) VALUES ('Chest', 'Red', '10.6', '50');
INSERT INTO `redventures`.`widget` (`name`, `price`, `melts`, `inventory`) VALUES ('Fire Sword', '150.75', 1, '15');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `inventory`) VALUES ('Elder Scrool', 'Blue', '2562.25', '2');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `melts`, `inventory`) VALUES ('Ice Staff', 'White', '541.65', 1, '54');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `inventory`) VALUES ('Duid Mask', 'Green', '12.65', '38');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `inventory`) VALUES ('Phoenix Pheather', 'Yellow', '1000.00', '2');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `melts`, `inventory`) VALUES ('Ghost Lamp', 'Dark Blue', '52.94', 1, '5');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `melts`, `inventory`) VALUES ('Key of Oracle', 'Gold', '12.54', 1, '1');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `inventory`) VALUES ('Thunder Bow', 'Brown', '65.54', '15');
INSERT INTO `redventures`.`widget` (`name`, `color`, `price`, `melts`, `inventory`) VALUES ('Amulet of Wisdom', 'Silver', '54.2', 1, '4');
