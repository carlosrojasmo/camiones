package main 

import (
	"fmt"
	"os"
	"sync"
	"time"
	"math/rand"
	"strconv"
)

var wg = &sync.WaitGroup{}

type paquete struct {
	idPaquete string
	tipo string
	valor int
	seguimiento int
	intentos int
	estado string
}

func newPaquete(idPaquete string, tipo string, valor int) *paquete{
	paqueteNuevo := paquete{idPaquete: idPaquete, tipo: tipo,valor: valor}
	random := rand.NewSource(time.Now().UnixNano())
	paqueteNuevo.seguimiento=(rand.New(random)).Intn(492829)
	paqueteNuevo.intentos= 0;
	paqueteNuevo.estado="en bodega"
	return &paqueteNuevo
}
type orden struct {
	timestamp time.Time
	idPaquete string
	tipo string
	nombre string
	valor int
	origen string
	destino string
	seguimiento int
}

func newOrden(tipo string,nombre string,valor int,origen string, destino string,id string) *orden{
	ordenNueva := orden{nombre: nombre,tipo: tipo,valor: valor ,origen: origen,idPaquete: id, destino: destino}
	random := rand.NewSource(time.Now().UnixNano())
	ordenNueva.seguimiento=(rand.New(random)).Intn(492829)
	ordenNueva.timestamp= time.Now()
	return &ordenNueva
}

type registroPorCamion struct {
	idPaquete string
	tipo string
	valor int
	origen string
	destino string
	intentos int
	entrega time.Time 

}


type camion struct{
	tipo string
	estado string
	carga [2]paquete
}

func newCamion(tipo string) *camion{
	camionNuevo := camion{tipo : tipo,estado : "Central",carga : [2]paquete{}}
	return &camionNuevo

}

func camionLaborando(tipo string){
	camionp := newCamion(tipo)
	cargaLenght := 0
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
	fmt.Println(camionp.tipo)
	for {
		if camionp.estado == "Central" {
			//pedir
			if true { // si salio bien el pedir
				camionp.carga[0] =  *newPaquete("idPaquete string", "tipo string", 10)
				cargaLenght = cargaLenght + 1
				if cargaLenght == 2 {
					camionp.estado = "Reparto"
					//elegir orden
				}
				fmt.Println(camionp.carga[0])
			}
		} else {
			//elegir paquete con mayor digni
			//aumentar intento en 1
			recibido = r.Intn(100)
			if recibido < 80 {
				cargaLenght = cargaLenght - 1
				if cargaLenght == 0 {
					camionp.estado = "Central"
				}


		     } else {
		     	//cambiar a otro paquete
		     }
		
	}
}
 func main() {

 	waitFor2 := 30

	if len(os.Args) > 1 {
		wait,err := strconv.Atoi(os.Args[1])
		if err != nil {
			wait = 30
		}
		waitFor2 = wait
		camionLaborando("normal")
	} 
    
	fmt.Println(waitFor2)
}