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
	address  = "10.10.28.10:50051"
)

var wg = &sync.WaitGroup{}


type entradaRegistroPorCamion struct {
	idPaquete string
	tipo string
	valor int64
	origen string
	destino string
	intentos int
	entrega time.Time 

}

func newEntrada(idPaquete string, tipo string, valor int64,origen string,destino string) *entradaRegistroPorCamion{
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
	fmt.Println("1")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
    fmt.Println("2")
	defer conn.Close()
	c := pb.NewOrdenServiceClient(conn)
    fmt.Println(3)
    //ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	camionp := newCamion(tipo)
	startwait := time.Now()
	var first int64
	first = 1
	cargaEntregada := []int{}
	cargaNoEntregada := []int{}
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
	fmt.Println("4")
	for {
		fmt.Println("camionp.tipo")
		if camionp.estado == "Central" {
			//pedir
			time.Sleep(4 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.GetPack(ctx, &pb.AskForPack{Tipo : camionp.tipo})
	        if err != nil {
	        	log.Fatalf("could not greet: %v", err)
	        }
	        fmt.Println(camionp.regi)
			if err == nil && r.GetIdPaquete() != "400" && r.GetIdPaquete() != ""{ // si salio bien el pedir
				camionp.regi = append(camionp.regi, *newEntrada(r.IdPaquete, r.Tipo, r.Valor,r.Origen,r.Destino))
				camionp.carga = append(camionp.carga,len(camionp.regi)-1)
				fmt.Println(len(camionp.regi)-1)
				if camionp.cargaLenght == 0{
					startwait = time.Now()
				}
				camionp.cargaLenght = camionp.cargaLenght + 1
				if camionp.cargaLenght == 2 {
					camionp.estado = "Reparto"
					first = 1
					//elegir orden
				} 
				fmt.Println(camionp.regi)
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
			paqueteEnEntrega := 0
			if camionp.cargaLenght > 1 {

				if camionp.regi[camionp.carga[0]].valor * first < camionp.regi[camionp.carga[1]].valor * first {
					paqueteEnEntrega = 1

				}
			}
			camionp.regi[camionp.carga[paqueteEnEntrega]].intentos = camionp.regi[camionp.carga[paqueteEnEntrega]].intentos + 1
			time.Sleep(10 * time.Second)
			recibido := r.Intn(100)
			fmt.Println(recibido)
			if recibido < 80 {
				camionp.regi[camionp.carga[paqueteEnEntrega]].entrega = time.Now()
				camionp.cargaLenght = camionp.cargaLenght - 1
				cargaEntregada = append(cargaEntregada,camionp.carga[paqueteEnEntrega])
				camionp.carga[paqueteEnEntrega] = camionp.carga[len(camionp.carga)-1] 
				camionp.carga[len(camionp.carga)-1] = 0   
				camionp.carga = camionp.carga[:len(camionp.carga)-1] 
		     } else {
		     	tipoPaq := camionp.regi[camionp.carga[paqueteEnEntrega]].tipo
		     	intentPaq := camionp.regi[camionp.carga[paqueteEnEntrega]].intentos
		     	if (tipoPaq == "retail" && intentPaq >= 3) || (tipoPaq != "retail" && intentPaq >= 2){
		     		camionp.cargaLenght = camionp.cargaLenght - 1
		     		cargaNoEntregada = append(cargaNoEntregada,camionp.carga[paqueteEnEntrega])
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
		     		fmt.Println("justo antes de entrega")
		     		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			        defer cancel()
		     		r, err := c.Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaEntregada[i]].idPaquete,
		     		 Entregado : true , Intentos : int64(camionp.regi[cargaEntregada[i]].intentos)})
		     		fmt.Println(r)
		     		if err != nil {
		     			log.Fatalf("could not greet: %v", err)
		     		}
		     	}
		     	for j := 0; j < len(cargaNoEntregada); j++ {
		     		time.Sleep(2 * time.Second)
		     		fmt.Println("justo antes de No entrega")
		     		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			        defer cancel()
		     		r, err := c.Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaNoEntregada[j]].idPaquete,
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

 	fmt.Println(waitFor2)

	if len(os.Args) > 1 {
		wait,err := strconv.Atoi(os.Args[1])
		if err != nil {
			wait = 30
		}
		waitFor2 = wait
	}

	fmt.Println(waitFor2)

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