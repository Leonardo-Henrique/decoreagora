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
  -- Table `decoreagora`.`users`
  -- -----------------------------------------------------
  CREATE TABLE IF NOT EXISTS `decoreagora`.`users` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(50) NOT NULL,
    `last_login` DATETIME NOT NULL,
    `public_id` VARCHAR(36) NOT NULL,
    PRIMARY KEY (`id`))
  ENGINE = InnoDB;


  -- -----------------------------------------------------
  -- Table `decoreagora`.`access_codes`
  -- -----------------------------------------------------
  CREATE TABLE IF NOT EXISTS `decoreagora`.`access_codes` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(255) NOT NULL,
    `used` TINYINT NOT NULL,
    `users_id` INT NOT NULL,
    `expire_at` DATETIME NULL COMMENT 'Timestamp',
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
    `original_image_key` VARCHAR(100) NOT NULL,
    `generated_image_key` VARCHAR(100),
    `prompt_description` VARCHAR(2000),
    `users_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_generated_images_users1_idx` (`users_id` ASC) VISIBLE,
    CONSTRAINT `fk_generated_images_users1`
      FOREIGN KEY (`users_id`)
      REFERENCES `decoreagora`.`users` (`id`)
      ON DELETE NO ACTION
      ON UPDATE NO ACTION)
  ENGINE = InnoDB;


  -- -----------------------------------------------------
  -- Table `decoreagora`.`available_credits`
  -- -----------------------------------------------------
  CREATE TABLE IF NOT EXISTS `decoreagora`.`available_credits` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `total` INT NOT NULL,
    `users_id` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_available_credits_users1_idx` (`users_id` ASC) VISIBLE,
    CONSTRAINT `fk_available_credits_users1`
      FOREIGN KEY (`users_id`)
      REFERENCES `decoreagora`.`users` (`id`)
      ON DELETE NO ACTION
      ON UPDATE NO ACTION)
  ENGINE = InnoDB;


  -- -----------------------------------------------------
  -- Table `decoreagora`.`subscriptions`
  -- -----------------------------------------------------
  CREATE TABLE IF NOT EXISTS `decoreagora`.`subscriptions` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `stripe_costumer_id` VARCHAR(255),
    `stripe_subscription_id` VARCHAR(255),
    `stripe_price_id` VARCHAR(255),
    `is_active` TINYINT,
    `tier` VARCHAR(100) NOT NULL,
    `user_email`VARCHAR(100),
    `users_id` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_subscriptions_users1_idx` (`users_id` ASC) VISIBLE,
    CONSTRAINT `fk_subscriptions_users1`
      FOREIGN KEY (`users_id`)
      REFERENCES `decoreagora`.`users` (`id`)
      ON DELETE NO ACTION
      ON UPDATE NO ACTION)
  ENGINE = InnoDB;

  -- -----------------------------------------------------
  -- Table `decoreagora`.`payment_history`
  -- -----------------------------------------------------
  CREATE TABLE IF NOT EXISTS `decoreagora`.`payment_history` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `stripe_payment_id` VARCHAR(255) NOT NULL,
    `public_id` VARCHAR(36) NOT NULL, 
    `stripe_customer_id`VARCHAR(255) NOT NULL,
    `processed_at` DATETIME NOT NULL,
    `stripe_price_id` VARCHAR(255) NOT NULL,
    `amount_paid` INT NOT NULL,
    `credits_received` INT NOT NULL,
    PRIMARY KEY (`id`))
  ENGINE = InnoDB;


  SET SQL_MODE=@OLD_SQL_MODE;
  SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
  SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
