-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema decoreagora
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema decoreagora
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `decoreagora` DEFAULT CHARACTER SET utf8 ;
USE `decoreagora` ;

-- -----------------------------------------------------
-- Table `decoreagora`.`available_credits`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `decoreagora`.`available_credits` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `total` INT NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `decoreagora`.`subscriptions`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `decoreagora`.`subscriptions` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `stripe_costumer_id` VARCHAR(255) NOT NULL,
  `stripe_subscription_id` VARCHAR(255) NOT NULL,
  `stripe_price_id` VARCHAR(255) NOT NULL,
  `is_active` TINYINT NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `decoreagora`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `decoreagora`.`users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `last_login` DATETIME NOT NULL,
  `public_id` VARCHAR(36) NOT NULL,
  `available_credits_id` INT NOT NULL,
  `subscriptions_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_users_available_credits1_idx` (`available_credits_id` ASC) VISIBLE,
  INDEX `fk_users_subscriptions1_idx` (`subscriptions_id` ASC) VISIBLE,
  CONSTRAINT `fk_users_available_credits1`
    FOREIGN KEY (`available_credits_id`)
    REFERENCES `decoreagora`.`available_credits` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_users_subscriptions1`
    FOREIGN KEY (`subscriptions_id`)
    REFERENCES `decoreagora`.`subscriptions` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `decoreagora`.`access_codes`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `decoreagora`.`access_codes` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(255) NOT NULL,
  `used` TINYINT NOT NULL,
  `users_id` INT NOT NULL,
  `expire_at` VARCHAR(45) NULL COMMENT 'Timestamp',
  PRIMARY KEY (`id`),
  INDEX `fk_access_codes_users_idx` (`users_id` ASC) VISIBLE,
  CONSTRAINT `fk_access_codes_users`
    FOREIGN KEY (`users_id`)
    REFERENCES `decoreagora`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `decoreagora`.`generated_images`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `decoreagora`.`generated_images` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `public_id` VARCHAR(36) NOT NULL,
  `url` VARCHAR(45) NOT NULL,
  `users_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_generated_images_users1_idx` (`users_id` ASC) VISIBLE,
  CONSTRAINT `fk_generated_images_users1`
    FOREIGN KEY (`users_id`)
    REFERENCES `decoreagora`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
