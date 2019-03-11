-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema site
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema site
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `site` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ;
USE `site` ;

-- -----------------------------------------------------
-- Table `site`.`comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `site`.`comments` (
  `id` MEDIUMTEXT NULL DEFAULT NULL,
  `postID` MEDIUMTEXT NULL DEFAULT NULL,
  `userID` MEDIUMTEXT NULL DEFAULT NULL,
  `author` MEDIUMTEXT NULL DEFAULT NULL,
  `body` MEDIUMTEXT NULL DEFAULT NULL,
  `created` DATETIME NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `site`.`posts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `site`.`posts` (
  `id` LONGTEXT NULL DEFAULT NULL,
  `userID` LONGTEXT NULL DEFAULT NULL,
  `author` LONGTEXT NULL DEFAULT NULL,
  `title` LONGTEXT NULL DEFAULT NULL,
  `body` LONGTEXT NULL DEFAULT NULL,
  `gamerRage` TINYINT(1) NULL DEFAULT NULL,
  `votes` INT(11) NULL DEFAULT NULL,
  `created` DATETIME NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `site`.`replies`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `site`.`replies` (
  `id` MEDIUMTEXT NULL DEFAULT NULL,
  `parentID` MEDIUMTEXT NULL DEFAULT NULL,
  `userID` MEDIUMTEXT NULL DEFAULT NULL,
  `author` MEDIUMTEXT NULL DEFAULT NULL,
  `body` MEDIUMTEXT NULL DEFAULT NULL,
  `created` DATETIME NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `site`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `site`.`users` (
  `id` MEDIUMTEXT NULL DEFAULT NULL,
  `username` MEDIUMTEXT NULL DEFAULT NULL,
  `password` MEDIUMTEXT NULL DEFAULT NULL,
  `created` DATETIME NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;


-- -----------------------------------------------------
-- Table `site`.`votes`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `site`.`votes` (
  `userID` TEXT NULL DEFAULT NULL,
  `postID` TEXT NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
