CREATE DATABASE IF NOT EXISTS pruebadb;
use pruebadb;
CREATE TABLE IF NOT EXISTS mensajeria(
	idmensajeria BIGINT NOT NULL AUTO_INCREMENT,
	usuarios VARCHAR(100) NOT NULL,
	mensajes LONGTEXT NOT NULL,
	fotos BLOB NULL,
	PRIMARY KEY(idmensajeria)
);
