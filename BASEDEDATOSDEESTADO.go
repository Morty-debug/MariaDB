package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"fmt"
	"database/sql"
    "log"
    "encoding/base64"
    _ "github.com/go-sql-driver/mysql"
)

type Respuesta struct {
	Profesion   			string
	PrimerNombre			string
	SegundoNombre			string
	TercerNombre			string
	PrimerApellido			string
	SegundoApellido			string
	Genero					string
	ColorPiel				string
	ColorOjos				string 
	ColorCabello			string
	Peso					string
	Estatura				string
	FechaNacimiento			string
	DUI 					string 
	Pasaporte 				string
	NumeroPartida 			string 
	FolioPartida 			string
	TomoPartida 			string
	TipoLibroPartida 		string 
	LibroPartida 			string
	AñoPartida 				string
	Foto 					string
}

var conexion = "dgmeconsultaafis:Dgme2019@tcp(10.10.1.160:3306)/afisgs"

func hasher(s string) []byte {
    val := sha256.Sum256([]byte(s))
    return val[:]
}

func authHandler(handler http.HandlerFunc, userhash, passhash []byte, realm string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || subtle.ConstantTimeCompare(hasher(user),
            userhash) != 1 || subtle.ConstantTimeCompare(hasher(pass), passhash) != 1 {
            w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
            http.Error(w, "No autorizado.", http.StatusUnauthorized)
            return
        }
        handler(w, r)
    }
}

func main() {
	//http.HandleFunc("/buscarPersona_xPasaporte", buscarPersona_xPasaporte)
	//http.ListenAndServe(":5002", nil)
	userhash := hasher("usu_cansilleria")
	passhash := hasher("emVsZGFfb2NhcmluYV9vZnRpbWU")
	realm := "BasicAuth necesita credenciales"
	http.HandleFunc("/buscarPersona_xPasaporte", authHandler(buscarPersona_xPasaporte, userhash, passhash, realm))
	http.HandleFunc("/buscarPersona_xDui", authHandler(buscarPersona_xDui, userhash, passhash, realm))
	http.HandleFunc("/buscarPersona_xPartida", authHandler(buscarPersona_xPartida, userhash, passhash, realm))
	http.ListenAndServe(":5002", nil)
}


func buscarPersona_xPasaporte(w http.ResponseWriter, r *http.Request) {
	pasaporte := r.FormValue("pasaporte")
	
	if len(pasaporte) <= 0 {
    	w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Formato no compatible\"}") 
		return
    }

	db, err := sql.Open("mysql", conexion)
    defer db.Close()

    if err != nil {
        log.Fatal(err) //error de conexion a base de datos
        return
    }
    
    query := `SELECT DISTINCT
				p.foto AS Foto,
				IF(pf.Profesion = '', '', pf.Profesion) AS Profesion,
				ps.FirstName AS PrimerNombre,
				IF(ps.MiddleName = '', '', ps.MiddleName) AS SegundoNombre,
				IF(ps.thirdName = '', '', ps.thirdName) AS TercerNombre,
				ps.LastName1 AS PrimerApellido,
				IFNULL((REPLACE(ps.LastName2, '|', '')),'') AS SegundoApellido,
				IF(ap.sexo = 0, 'M', 'F') AS Genero,
				pl.descripcion AS ColorPiel,
				pe.descripcion AS ColorCabello,
				o.descripcion AS ColorOjos,
				ap.estatura AS Estatura,
				ap.peso AS peso,
				ps.BirthDate AS FechaNacimiento,
				p.numpas AS NumeroPasaporte,
				ps.NID AS DUI,
				ap.numpar AS NumeroPartida,
				ap.foliopar AS FolioPartida,
				ap.idTomo AS TomoPartida,
				ap.Idtiplibpar AS TipoLibroPartida,
				ap.libpar AS LibroPartida,
				ap.anopar AS AnioPartida
			FROM afisgs.persons ps
			INNER JOIN pasaportes.pasaporte p ON ps.ID = p.idperson
			INNER JOIN pasaportes.anexopersonas ap ON ap.idPerson = p.idperson
			INNER JOIN pasaportes.piel pl ON pl.idpiel = ap.idpiel
			LEFT JOIN pasaportes.profesion pf ON pf.idProfesion = ap.idprofe
			INNER JOIN pasaportes.pelo pe ON pe.idpelo = ap.idpelo
			INNER JOIN pasaportes.ojos o ON o.idojos = ap.idojos
			INNER JOIN pasaportes.ciudades c ON c.idCiudades = ap.idciudadnac
			INNER JOIN pasaportes.estadoprovincia e ON e.idEstado = c.idEstado
			INNER JOIN pasaportes.pais pais ON pais.idPais = e.idPais
			INNER JOIN pasaportes.monitor moni ON moni.Idpasaporte = p.idpasaporte
			LEFT JOIN pasaportes.casoespecial cas ON cas.idCasoEspecial = ap.idCasoEspecial
			INNER JOIN pasaportes.delegacion del ON moni.iddelegacion = del.idDelegacion 
			LEFT OUTER JOIN pasaportes.revalidacion rev ON rev.idpasaporte = p.idpasaporte 
			WHERE p.idpaso IN (4,6) 
			AND p.numlib != '' 
			AND p.numlib IS NOT NULL
			AND p.numpas = '`+pasaporte+`';`

    rows, err := db.Query(query)
    if err != nil {
        panic(err.Error()) 
        return
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) 
        return
    }

    values := make([]sql.RawBytes, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    if rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        Estructura := Respuesta{string(values[1]), string(values[2]), string(values[3]), string(values[4]), string(values[5]), string(values[6]), string(values[7]), string(values[8]), string(values[10]), string(values[9]), string(values[12]), string(values[11]), string(values[13]), string(values[15]), string(values[14]), string(values[16]), string(values[17]), string(values[18]), string(values[19]), string(values[20]), string(values[21]), base64.StdEncoding.EncodeToString(values[0]) }
        
        js, err := json.Marshal(Estructura)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"Estado\":\"No se logro crear JSON\"}") 
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
    } else {
    	Estructura := Respuesta{ }
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
}




func buscarPersona_xDui(w http.ResponseWriter, r *http.Request) {
	dui := r.FormValue("dui")
	
	if len(dui) <= 0 {
    	w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Formato no compatible\"}") 
		return
    }

	db, err := sql.Open("mysql", conexion)
    defer db.Close()

    if err != nil {
        log.Fatal(err) //error de conexion a base de datos
        return
    }
    
    query := `SELECT DISTINCT
				p.foto AS Foto,
				IF(pf.Profesion = '', '', pf.Profesion) AS Profesion,
				ps.FirstName AS PrimerNombre,
				IF(ps.MiddleName = '', '', ps.MiddleName) AS SegundoNombre,
				IF(ps.thirdName = '', '', ps.thirdName) AS TercerNombre,
				ps.LastName1 AS PrimerApellido,
				IFNULL((REPLACE(ps.LastName2, '|', '')),'') AS SegundoApellido,
				IF(ap.sexo = 0, 'M', 'F') AS Genero,
				pl.descripcion AS ColorPiel,
				pe.descripcion AS ColorCabello,
				o.descripcion AS ColorOjos,
				ap.estatura AS Estatura,
				ap.peso AS peso,
				ps.BirthDate AS FechaNacimiento,
				p.numpas AS NumeroPasaporte,
				ps.NID AS DUI,
				ap.numpar AS NumeroPartida,
				ap.foliopar AS FolioPartida,
				ap.idTomo AS TomoPartida,
				ap.Idtiplibpar AS TipoLibroPartida,
				ap.libpar AS LibroPartida,
				ap.anopar AS AnioPartida
			FROM afisgs.persons ps
			INNER JOIN pasaportes.pasaporte p ON ps.ID = p.idperson
			INNER JOIN pasaportes.anexopersonas ap ON ap.idPerson = p.idperson
			INNER JOIN pasaportes.piel pl ON pl.idpiel = ap.idpiel
			LEFT JOIN pasaportes.profesion pf ON pf.idProfesion = ap.idprofe
			INNER JOIN pasaportes.pelo pe ON pe.idpelo = ap.idpelo
			INNER JOIN pasaportes.ojos o ON o.idojos = ap.idojos
			INNER JOIN pasaportes.ciudades c ON c.idCiudades = ap.idciudadnac
			INNER JOIN pasaportes.estadoprovincia e ON e.idEstado = c.idEstado
			INNER JOIN pasaportes.pais pais ON pais.idPais = e.idPais
			INNER JOIN pasaportes.monitor moni ON moni.Idpasaporte = p.idpasaporte
			LEFT JOIN pasaportes.casoespecial cas ON cas.idCasoEspecial = ap.idCasoEspecial
			INNER JOIN pasaportes.delegacion del ON moni.iddelegacion = del.idDelegacion 
			LEFT OUTER JOIN pasaportes.revalidacion rev ON rev.idpasaporte = p.idpasaporte 
			WHERE p.idpaso IN (4,6) 
			AND p.numlib != '' 
			AND p.numlib IS NOT NULL
			AND ps.NID = '`+dui+`';`

    rows, err := db.Query(query)
    if err != nil {
        panic(err.Error()) 
        return
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) 
        return
    }

    values := make([]sql.RawBytes, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    if rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        Estructura := Respuesta{string(values[1]), string(values[2]), string(values[3]), string(values[4]), string(values[5]), string(values[6]), string(values[7]), string(values[8]), string(values[10]), string(values[9]), string(values[12]), string(values[11]), string(values[13]), string(values[15]), string(values[14]), string(values[16]), string(values[17]), string(values[18]), string(values[19]), string(values[20]), string(values[21]), base64.StdEncoding.EncodeToString(values[0]) }
        
        js, err := json.Marshal(Estructura)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"Estado\":\"No se logro crear JSON\"}") 
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
    } else {
    	Estructura := Respuesta{ }
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
}




func buscarPersona_xPartida(w http.ResponseWriter, r *http.Request) {
	numeropartida := r.FormValue("numeropartida")
	tipolibropartida := r.FormValue("tipolibropartida")
	libropartida := r.FormValue("libropartida")
	foliopartida := r.FormValue("foliopartida")
	tomopartida := r.FormValue("tomopartida")
	aniopartida := r.FormValue("añopartida")
	
	if len(numeropartida) <= 0 || len(tipolibropartida) <= 0 || len(libropartida) <= 0 || len(foliopartida) <= 0 || len(tomopartida) <= 0 || len(aniopartida) <= 0 {
    	w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Formato no compatible\"}") 
		return
    }

	db, err := sql.Open("mysql", conexion)
    defer db.Close()

    if err != nil {
        log.Fatal(err) //error de conexion a base de datos
        return
    }
    
    query := `SELECT DISTINCT
				p.foto AS Foto,
				IF(pf.Profesion = '', '', pf.Profesion) AS Profesion,
				ps.FirstName AS PrimerNombre,
				IF(ps.MiddleName = '', '', ps.MiddleName) AS SegundoNombre,
				IF(ps.thirdName = '', '', ps.thirdName) AS TercerNombre,
				ps.LastName1 AS PrimerApellido,
				IFNULL((REPLACE(ps.LastName2, '|', '')),'') AS SegundoApellido,
				IF(ap.sexo = 0, 'M', 'F') AS Genero,
				pl.descripcion AS ColorPiel,
				pe.descripcion AS ColorCabello,
				o.descripcion AS ColorOjos,
				ap.estatura AS Estatura,
				ap.peso AS peso,
				ps.BirthDate AS FechaNacimiento,
				p.numpas AS NumeroPasaporte,
				ps.NID AS DUI,
				ap.numpar AS NumeroPartida,
				ap.foliopar AS FolioPartida,
				ap.idTomo AS TomoPartida,
				ap.Idtiplibpar AS TipoLibroPartida,
				ap.libpar AS LibroPartida,
				ap.anopar AS AnioPartida
			FROM afisgs.persons ps
			INNER JOIN pasaportes.pasaporte p ON ps.ID = p.idperson
			INNER JOIN pasaportes.anexopersonas ap ON ap.idPerson = p.idperson
			INNER JOIN pasaportes.piel pl ON pl.idpiel = ap.idpiel
			LEFT JOIN pasaportes.profesion pf ON pf.idProfesion = ap.idprofe
			INNER JOIN pasaportes.pelo pe ON pe.idpelo = ap.idpelo
			INNER JOIN pasaportes.ojos o ON o.idojos = ap.idojos
			INNER JOIN pasaportes.ciudades c ON c.idCiudades = ap.idciudadnac
			INNER JOIN pasaportes.estadoprovincia e ON e.idEstado = c.idEstado
			INNER JOIN pasaportes.pais pais ON pais.idPais = e.idPais
			INNER JOIN pasaportes.monitor moni ON moni.Idpasaporte = p.idpasaporte
			LEFT JOIN pasaportes.casoespecial cas ON cas.idCasoEspecial = ap.idCasoEspecial
			INNER JOIN pasaportes.delegacion del ON moni.iddelegacion = del.idDelegacion 
			LEFT OUTER JOIN pasaportes.revalidacion rev ON rev.idpasaporte = p.idpasaporte 
			WHERE p.idpaso IN (4,6) 
			AND p.numlib != '' 
			AND p.numlib IS NOT NULL
			AND ap.numpar = '`+numeropartida+`' AND ap.foliopar = '`+foliopartida+`' AND ap.idTomo = '`+tomopartida+`' AND ap.Idtiplibpar = '`+tipolibropartida+`' AND ap.libpar = '`+libropartida+`' AND ap.anopar = '`+aniopartida+`';`

    rows, err := db.Query(query)
    if err != nil {
        panic(err.Error()) 
        return
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) 
        return
    }

    values := make([]sql.RawBytes, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    if rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        Estructura := Respuesta{string(values[1]), string(values[2]), string(values[3]), string(values[4]), string(values[5]), string(values[6]), string(values[7]), string(values[8]), string(values[10]), string(values[9]), string(values[12]), string(values[11]), string(values[13]), string(values[15]), string(values[14]), string(values[16]), string(values[17]), string(values[18]), string(values[19]), string(values[20]), string(values[21]), base64.StdEncoding.EncodeToString(values[0]) }
        
        js, err := json.Marshal(Estructura)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"Estado\":\"No se logro crear JSON\"}") 
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
    } else {
    	Estructura := Respuesta{ }
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
}

