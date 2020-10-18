package main 

import (
    "context"
	"fmt"
	"os"
	"sync"
	"time"
	"math/rand"
	"strconv"
	"log"
	"google.golang.org/grpc"
	pb "../proto"
)

const (
	address  = "localhost:50054"
)

var wg = &sync.WaitGroup{}


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
	defer wg.Done()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewCamionDeliveryClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	camionp := newCamion(tipo)
	startwait := time.Now()
	first := 1
	cargaEntregada := []int{}
	cargaNoEntregada := []int{}
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
	fmt.Println(camionp.tipo)
	for {
		fmt.Println("camionp.tipo")
		if camionp.estado == "Central" {
			//pedir
			time.Sleep(4 * time.Second)
			r, err := c. GetPack(ctx, &pb.AskForPack{Tipo : camionp.tipo})
	        if err != nil {
	        	fmt.Println("could not greet")
	        }
			if err == nil && r.GetIdPaquete() != "400" { // si salio bien el pedir
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
				} 
				fmt.Println(camionp.regi[0])
			}
			if camionp.cargaLenght == 1 {
				dif :=  time.Now().Sub(startwait)
				if dif.Seconds() > waitFor2  {
					camionp.estado = "Reparto"
					first = 1
				}
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
			time.Sleep(10 * time.Second)
			recibido := r.Intn(100)
			fmt.Println(recibido)
			if recibido < 80 {
				camionp.cargaLenght = camionp.cargaLenght - 1
				cargaEntregada = append(cargaEntregada,paqueteEnEntrega)
				camionp.carga[paqueteEnEntrega] = camionp.carga[len(camionp.carga)-1] 
				camionp.carga[len(camionp.carga)-1] = 0   
				camionp.carga = camionp.carga[:len(camionp.carga)-1] 
		     } else {
		     	tipoPaq := camionp.regi[paqueteEnEntrega].tipo
		     	intentPaq := camionp.regi[paqueteEnEntrega].intentos
		     	if (tipoPaq == "retail" && intentPaq >= 3) || (tipoPaq != "retail" && intentPaq >= 2){
		     		camionp.cargaLenght = camionp.cargaLenght - 1
		     		cargaNoEntregada = append(cargaNoEntregada,paqueteEnEntrega)
		     		camionp.carga[paqueteEnEntrega] = camionp.carga[len(camionp.carga)-1] 
				    camionp.carga[len(camionp.carga)-1] = 0   
				    camionp.carga = camionp.carga[:len(camionp.carga)-1] 
		     	} 
		     	first = -first
		     }
		     if camionp.cargaLenght == 0 {
		     	camionp.estado = "Central"
		     	for i := 0; i < len(cargaEntregada); i++ {
		     		time.Sleep(2 * time.Second)
		     		r, err := c. Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaEntregada[i]].idPaquete,
		     		 Entregado : true , Intentos : int64(camionp.regi[cargaEntregada[i]].intentos)})
		     		fmt.Println(r)
		     		if err != nil {
		     			log.Fatalf("could not greet: %v", err)
		     		}
		     	}
		     	for j := 0; j < len(cargaNoEntregada); j++ {
		     		time.Sleep(2 * time.Second)
		     		r, err := c. Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaNoEntregada[j]].idPaquete,
		     		 Entregado : false ,Intentos : int64(camionp.regi[cargaNoEntregada[j]].intentos)})
		     		fmt.Println(r)
		     		if err != nil {
		     			log.Fatalf("could not greet: %v", err)
		     		}
		    	}
		    	cargaEntregada = []int{}
	            cargaNoEntregada = []int{}
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
	}

	wg.Add(1)
	go camionLaborando("normal",float64(waitFor2)) 

	time.Sleep(6 * time.Second)

    wg.Add(1)
	go camionLaborando("retail",float64(waitFor2))

	time.Sleep(6 * time.Second)

    wg.Add(1)
	go camionLaborando("retail",float64(waitFor2))

    wg.Wait()
	fmt.Println(waitFor2)
}