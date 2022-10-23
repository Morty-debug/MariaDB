package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"database/sql"
	"encoding/base64"
	_ "github.com/go-sql-driver/mysql"
)

var conexion = "root:123456@tcp(127.0.0.1:3306)/pruebadb"

func main() {
	var opcion string
	var ejecutar = true
	for ejecutar {
		fmt.Printf("1 Insertar Datos\n")
		fmt.Printf("2 Mostrar Datos\n")
		fmt.Printf("3 Salir\n>")
		fmt.Scanf("%s", &opcion)
		screen.Clear()
		screen.MoveTopLeft()
		if opcion == "1" {
			insertar() //ejemplo de insert
		} else if opcion == "2" {
			mostrar() //ejemplo de select
		} else if opcion == "3" {
			ejecutar = false			
		}
	}
}

func insertar() {
	var usuario, mensaje string
	/*foto Byte[]*/
	foto, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAMAAAC6sdbXAAADAFBMVEUAAAD///8CAgIDAwMEBAQFBQUGBgYHBwcICAgJCQkKCgoLCwsMDAwNDQ0ODg4PDw8QEBARERESEhITExMUFBQVFRUWFhYXFxcYGBgZGRkaGhobGxscHBwdHR0eHh4fHx8gICAhISEiIiIjIyMkJCQlJSUmJiYnJycoKCgpKSkqKiorKyssLCwtLS0uLi4vLy8wMDAxMTEyMjIzMzM0NDQ1NTU2NjY3Nzc4ODg5OTk6Ojo7Ozs8PDw9PT0+Pj4/Pz9AQEBBQUFCQkJDQ0NERERFRUVGRkZHR0dISEhJSUlKSkpLS0tMTExNTU1OTk5PT09QUFBRUVFSUlJTU1NUVFRVVVVWVlZXV1dYWFhZWVlaWlpbW1tcXFxdXV1eXl5fX19gYGBhYWFiYmJjY2NkZGRlZWVmZmZnZ2doaGhpaWlqampra2tsbGxtbW1ubm5vb29wcHBxcXFycnJzc3N0dHR1dXV2dnZ3d3d4eHh5eXl6enp7e3t8fHx9fX1+fn5/f3+AgICBgYGCgoKDg4OEhISFhYWGhoaHh4eIiIiJiYmKioqLi4uMjIyNjY2Ojo6Pj4+QkJCRkZGSkpKTk5OUlJSVlZWWlpaXl5eYmJiZmZmampqbm5ucnJydnZ2enp6fn5+goKChoaGioqKjo6OkpKSlpaWmpqanp6eoqKipqamqqqqrq6usrKytra2urq6vr6+wsLCxsbGysrKzs7O0tLS1tbW2tra3t7e4uLi5ubm6urq7u7u8vLy9vb2+vr6/v7/AwMDBwcHCwsLDw8PExMTFxcXGxsbHx8fIyMjJycnKysrLy8vMzMzNzc3Ozs7Pz8/Q0NDR0dHS0tLT09PU1NTV1dXW1tbX19fY2NjZ2dna2trb29vc3Nzd3d3e3t7f39/g4ODh4eHi4uLj4+Pk5OTl5eXm5ubn5+fo6Ojp6enq6urr6+vs7Ozt7e3u7u7v7+/w8PDx8fHy8vLz8/P09PT19fX29vb39/f4+Pj5+fn6+vr7+/v8/Pz9/f3+/v7////5JLncAAAAFElEQVQI12NgZECCIMDAACEYgBgAARoAEVcwRFUAAAAASUVORK5CYII=")
	/*ingresamos variable a insertar*/
	fmt.Printf("USUARIO: ")
	fmt.Scanf("%s", &usuario)
	fmt.Printf("MENSAJE: ")
	fmt.Scanf("%s", &mensaje)
	/*establecemos la conexion*/
	db, err := sql.Open("mysql", conexion)
	/*en caso de error detenemos todo*/
	if err != nil {
		fmt.Println(err.Error())
		db.Close()
		return
	}
	defer db.Close()
	/*ejecutamos insert*/
	_, err = db.Query("INSERT INTO mensajeria(usuarios,mensajes,fotos) VALUES(?, ?, ?)",usuario, mensaje, foto)
	/*en caso de error detenemos todo*/
	if err != nil {
		fmt.Println(err.Error())
		db.Close()
		return
	}
	db.Close()
}

func mostrar() {
	var usuario, mensaje string
	var foto []byte
	/*establecemos la conexion*/
	db, err := sql.Open("mysql", conexion)
	defer db.Close()
	/*en caso de error detenemos todo*/
	if err != nil {
		fmt.Println(err.Error())
		db.Close()
		return
	}
	/*ejecutamos select*/
	res, err := db.Query("SELECT usuarios, mensajes, fotos FROM mensajeria")
	defer res.Close()
	/*en caso de error detenemos todo*/
	if err != nil {
		fmt.Println(err.Error())
		db.Close()
		res.Close()
		return
	}
	/*recorremos filas*/
	for res.Next() {
		/*capturamos cada dato de la fila*/
		err := res.Scan(&usuario, &mensaje, &foto)
		/*en caso de error detenemos todo*/
		if err != nil {
			fmt.Println(err.Error())
			db.Close()
			res.Close()
			return
		}
		/*mostramos en terminal*/
		fmt.Printf("USUARIO=%s MENSAJE=%s FOTO=%s\n\n", usuario, mensaje, base64.StdEncoding.EncodeToString(foto))
	}
	/*generamos una pausa*/
	fmt.Printf("PRECIONE UNA TECLA PARA CONTINUAR...\n")
	fmt.Scanf("%s")
	db.Close()
}
