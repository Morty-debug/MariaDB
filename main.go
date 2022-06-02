//https://www.geeksforgeeks.org/how-to-use-go-with-mysql/

//https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://zetcode.com/golang/mysql/
//https://github.com/go-sql-driver/mysql/wiki/Examples


package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

type Recepcion struct {
	ID_Sucursal			string `json:"id_sucursal"`
}

type Respuesta struct {
	Sucursal   			string
	Citas				Disponibilidades
}

type Disponibilidades struct {
	Disponible []Disponibilidad
}

type Disponibilidad struct {
	Disponibles			string
	Fecha				string
	Horario				string
}

func (box *Disponibilidades) AddItem(item Disponibilidad) {
	box.Disponible = append(box.Disponible, item)
}

var conexion = "ricardo.valladares:Ricardo2019@tcp(192.100.1.240:3306)/citaspasaporte"

/*Ejemplo de JSON a RETORNAR:
{
	"Sucursal":"Cascadas",
	"Citas": {
        "Disponible": [
        {
			"Disponibles": "0",
			"Fecha": "2022-04-30",
			"Horario": "09:00 - 12:00"
        },{
			"Disponibles": "150",
			"Fecha": "2022-04-30",
			"Horario": "09:00 - 12:00"
		}
	]
}
*/

/*Ejemplo de JSON a RECIBIR:
{
	"ID_Sucursal":"1"
}
*/




func main() {
	http.HandleFunc("/wscitas", wscitas)
	http.ListenAndServe(":8001", nil)
}


func wscitas(w http.ResponseWriter, r *http.Request) {
	
	box := Disponibilidades{}
	var sucursal string 

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Formato no compatible\"}") 
		return
	}

	var data Recepcion
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Estructura no compatible\"}") 
		return
	}
	
	fmt.Println("Sucursal: ", data.ID_Sucursal)

	db, err := sql.Open("mysql", conexion)
    defer db.Close()

    if err != nil {
        log.Fatal(err) //error de conexion a base de datos
    }
    
    query := `SELECT 
			(SELECT sucursal FROM sucursal WHERE idSucursal=`+data.ID_Sucursal+`) AS Sucursal,
			tabla.Disponibles, 
			tabla.Fechas, 
			tabla.HoraInicio, 
			tabla.HoraFin,
			CONCAT('[',(SELECT GROUP_CONCAT('"',che.hora,'"') FROM citaspasaporte.horas_espera che WHERE che.hora BETWEEN TIME(tabla.HoraInicio) AND TIME(tabla.HoraFin) ORDER BY che.hora ASC),']') AS Horas 
			FROM (SELECT 
				IF( ((SELECT limite FROM citaspasaporte.horariosucursal h2 WHERE h2.idSucursal = `+data.ID_Sucursal+` AND h2.idTipoDia=1 LIMIT 1) - (SELECT COUNT(*) FROM citaspasaporte.detallecita d2 WHERE d2.idSucursal=`+data.ID_Sucursal+` AND d2.idFecha=f.idFecha LIMIT 1))<=0, 0, ((SELECT limite FROM citaspasaporte.horariosucursal h2 WHERE h2.idSucursal = `+data.ID_Sucursal+` AND h2.idTipoDia=1 LIMIT 1) - (SELECT COUNT(*) FROM citaspasaporte.detallecita d2 WHERE d2.idSucursal=`+data.ID_Sucursal+` AND d2.idFecha=f.idFecha LIMIT 1))  ) AS Disponibles,
				f.Fecha AS Fechas,
				(SELECT th1.horaInicio FROM citaspasaporte.horariosucursal h4 INNER JOIN citaspasaporte.tipohorario th1 ON th1.idTipoHorario = h4.idTipoHorario WHERE h4.idSucursal = `+data.ID_Sucursal+` AND h4.idTipoDia=1 LIMIT 1) AS HoraInicio,
				(SELECT th1.horaFin FROM citaspasaporte.horariosucursal h4 INNER JOIN citaspasaporte.tipohorario th1 ON th1.idTipoHorario = h4.idTipoHorario WHERE h4.idSucursal = `+data.ID_Sucursal+` AND h4.idTipoDia=1 LIMIT 1) AS HoraFin
				FROM citaspasaporte.fecha f WHERE f.Fecha>DATE(NOW()) AND f.estado=1 AND f.tipoDia=1	
				UNION ALL
				SELECT 
				IF( ((SELECT limite FROM citaspasaporte.horariosucursal h2 WHERE h2.idSucursal = `+data.ID_Sucursal+` AND h2.idTipoDia=2 LIMIT 1) - (SELECT COUNT(*) FROM citaspasaporte.detallecita d2 WHERE d2.idSucursal=`+data.ID_Sucursal+` AND d2.idFecha=f.idFecha LIMIT 1))<=0, 0, ((SELECT limite FROM citaspasaporte.horariosucursal h2 WHERE h2.idSucursal = `+data.ID_Sucursal+` AND h2.idTipoDia=2 LIMIT 1) - (SELECT COUNT(*) FROM citaspasaporte.detallecita d2 WHERE d2.idSucursal=`+data.ID_Sucursal+` AND d2.idFecha=f.idFecha LIMIT 1))  ) AS Disponibles,
				f.Fecha AS Fechas,
				(SELECT th1.horaInicio FROM citaspasaporte.horariosucursal h4 INNER JOIN citaspasaporte.tipohorario th1 ON th1.idTipoHorario = h4.idTipoHorario WHERE h4.idSucursal = `+data.ID_Sucursal+` AND h4.idTipoDia=2 LIMIT 1) AS HoraInicio,
				(SELECT th1.horaFin FROM citaspasaporte.horariosucursal h4 INNER JOIN citaspasaporte.tipohorario th1 ON th1.idTipoHorario = h4.idTipoHorario WHERE h4.idSucursal = `+data.ID_Sucursal+` AND h4.idTipoDia=2 LIMIT 1) AS HoraFin
				FROM citaspasaporte.fecha f WHERE f.Fecha>DATE(NOW()) AND f.estado=1 AND f.tipoDia=2
			) AS tabla
			ORDER BY tabla.Fechas;`

    rows, err := db.Query(query)
    if err != nil {
        panic(err.Error()) 
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) 
    }

    values := make([]sql.RawBytes, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }
        
        var disponibles string
        var fechas string
        var horainicio string
        var horafin string
        
        sucursal = string(values[0])
        disponibles = string(values[1])
        fechas = string(values[2])
        horainicio = string(values[3])
        horafin = string(values[4])

        item := Disponibilidad{disponibles, fechas, horainicio+" - "+horafin}

        if len(disponibles) > 0 {
        	box.AddItem(item)
        }
    }

    
    fmt.Println("Sucursal: ", sucursal)

	Estructura := Respuesta{sucursal, box}
	js, err := json.Marshal(Estructura)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"No se logro crear JSON\"}") 
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
