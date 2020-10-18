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

type entradaRegistroPorCamion struct {
	idPaquete string
	tipo string
	valor int
	origen string
	destino string
	intentos int
	entrega time.Time 

}

func newEntrada(idPaquete string, tipo string, valor int,origen string,destino string) *entradaRegistroPorCamion{
	entrada := entradaRegistroPorCamion{idPaquete: idPaquete, tipo: tipo,valor: valor,origen : origen,destino : destino}
	entrada.intentos = 0
	return &entrada
}

type camion struct{
	tipo string
	estado string
	regi []entradaRegistroPorCamion
	carga []int
	cargaLenght int
}

func newCamion(tipo string) *camion{
	camionNuevo := camion{tipo : tipo,estado : "Central",cargaLenght : 0}
	return &camionNuevo

}


func camionLaborando(tipo string,waitFor2 float64){
	camionp := newCamion(tipo)
	startwait := time.Now()
	first := 1
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
	fmt.Println(camionp.tipo)
	for {
		fmt.Println("camionp.tipo")
		if camionp.estado == "Central" {
			//pedir
			if true { // si salio bien el pedir
				camionp.regi = append(camionp.regi, *newEntrada("idPaquete string", "tipo string", 10,"micasa","tucasa"))
				camionp.carga = append(camionp.carga,len(camionp.regi)-1)
				if camionp.cargaLenght == 0{
					startwait = time.Now()
				}
				camionp.cargaLenght = camionp.cargaLenght + 1
				if camionp.cargaLenght == 2 {
					camionp.estado = "Reparto"
					first = 1
					//elegir orden
				} else if camionp.cargaLenght == 1 {
					dif :=  time.Now().Sub(startwait)
					if dif.Seconds() > waitFor2  {

					}
					
				}
				fmt.Println(camionp.regi[0])
			}
		} else {
			//elegir paquete con mayor digni
			paqueteEnEntrega := camionp.carga[0]
			if camionp.cargaLenght > 1 {

				if camionp.regi[camionp.carga[0]].valor * first < camionp.regi[camionp.carga[1]].valor * first {
					paqueteEnEntrega = camionp.carga[1]

				}
			}
			camionp.regi[paqueteEnEntrega].intentos = camionp.regi[paqueteEnEntrega].intentos + 1
			recibido := r.Intn(100)
			fmt.Println(recibido)
			if recibido < 80 {
				camionp.cargaLenght = camionp.cargaLenght - 1
				if camionp.cargaLenght == 0 {
					camionp.estado = "Central"
				}


		     } else {
		     	first = -first
		     }
		
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
		camionLaborando("normal",float64(waitFor2))
	} 
    
	fmt.Println(waitFor2)
}