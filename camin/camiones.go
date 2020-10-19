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

/*Esta funcion simula el funcionamiento de un camion ,recibe el tipo de camion 
,la epspera para recibir un segundo paquete .*/
func camionLaborando(tipo string,waitFor2 float64,waitForDeliverPack float64){
	defer wg.Done()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOrdenServiceClient(conn)
	camionp := newCamion(tipo) //Creamos nuestro camion
	startwait := time.Now() //Definimos starwait
	var first int64 // Esta variable nos ayudara a elegir que paquete entregamos primero
	first = 1
	cargaEntregada := []int{} //Estos 2 slices nos indicara que carga fue entregada y cual no
	cargaNoEntregada := []int{}
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
	for {
		if camionp.estado == "Central" {//si el camion esta en Central este pide paquetes periodicamente
			//pedir
			time.Sleep(4 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.GetPack(ctx, &pb.AskForPack{Tipo : camionp.tipo}) //Pedimos paquete
	        if err != nil {
	        	fmt.Println("Fallo de conexion en pedir paquete")
		        fmt.Println("could not greet: ", err)
	        }
	        
			if err == nil && r.GetIdPaquete() != "400" && r.GetIdPaquete() != ""{ // Si salio bien el pedir
				camionp.regi = append(camionp.regi, *newEntrada(r.IdPaquete, r.Tipo, r.Valor,r.Origen,r.Destino))
				camionp.carga = append(camionp.carga,len(camionp.regi)-1)//En el atributo regi guardamos la info del paquete
				                                                      //y en el atributo carga guardamos su indice en regi
				if camionp.cargaLenght == 0{ //Si fue el primer paquete en ingresar empesamos a esperar un segundo paquete
					startwait = time.Now()
				}
				camionp.cargaLenght = camionp.cargaLenght + 1 //Aumentamos el tamaño de la carga
				if camionp.cargaLenght == 2 { //si llego el 2 paquete el camion sale
					camionp.estado = "Reparto"
					first = 1
					
				} 
			}
			if camionp.cargaLenght == 1 { //Calculamos la diferencia en el tiempo que empezamos a esperar y 
				dif :=  time.Now().Sub(startwait) //el tiempo actual si es mayor a waitfor2 el camion sale
				if dif.Seconds() > waitFor2  {
					camionp.estado = "Reparto"
					first = 1
				}
			}
		} else {
			//Primero elegimos que paquete se entrega primero
			paqueteEnEntrega := 0
			if camionp.cargaLenght > 1 {//si el camion lleva un paquete no es necesario elegir
                /*Si el camion esta recien saliendo de la central se entrega el paquete con el mayor valor
                en caso contrario ,si fallamos un intento, la variable first cambia el sentido de la condicion 
                pasando a entregar el paquete de menor valor*/ 
				if camionp.regi[camionp.carga[0]].valor * first < camionp.regi[camionp.carga[1]].valor * first {
					paqueteEnEntrega = 1

				}
			}
			//Aumentamos en 1 la cantidad de intentos del paquete 
			camionp.regi[camionp.carga[paqueteEnEntrega]].intentos = camionp.regi[camionp.carga[paqueteEnEntrega]].intentos + 1
			time.Sleep(time.Duration(waitForDeliverPack) * time.Second) //Simulamos la espera en que el camion entrega el paquete
			recibido := r.Intn(100) //Con un random determinamos si este fue recibido
	
			if recibido < 80 { //El 80% de las veces es recibido
				camionp.regi[camionp.carga[paqueteEnEntrega]].entrega = time.Now() //Marcamos la fecha de entrega
				camionp.cargaLenght = camionp.cargaLenght - 1 //Lo descontamos de la carga
				cargaEntregada = append(cargaEntregada,camionp.carga[paqueteEnEntrega]) //Lo añadimos a cargaEntregada para 
				camionp.carga[paqueteEnEntrega] = camionp.carga[len(camionp.carga)-1]  //el posterior reporte
				camionp.carga[len(camionp.carga)-1] = 0   
				camionp.carga = camionp.carga[:len(camionp.carga)-1] //Las 3 lineas anteriores borran el paquete de la carga
		     } else { //Caso contrario 
		     	tipoPaq := camionp.regi[camionp.carga[paqueteEnEntrega]].tipo 
		     	intentPaq := camionp.regi[camionp.carga[paqueteEnEntrega]].intentos
		     	if (tipoPaq == "retail" && intentPaq >= 3) || (tipoPaq != "retail" && intentPaq >= 2){ //verificamos si supero
		     		camionp.cargaLenght = camionp.cargaLenght - 1                                //el numero de intentos limite
		     		cargaNoEntregada = append(cargaNoEntregada,camionp.carga[paqueteEnEntrega])
		     		camionp.carga[paqueteEnEntrega] = camionp.carga[len(camionp.carga)-1] 
				    camionp.carga[len(camionp.carga)-1] = 0   
				    camionp.carga = camionp.carga[:len(camionp.carga)-1]  //Eliminamos el elemento de la carga
		     	} 
		     	first = -first
		     }
		     if camionp.cargaLenght == 0 {// Cuando se acaba la carga el camion vuelve a la central
		     	camionp.estado = "Central"
		     	for i := 0; i < len(cargaEntregada); i++ {//Reportamos los paquetes entregados
		     		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			        defer cancel()
		     		_ , err := c.Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaEntregada[i]].idPaquete,
		     		 Entregado : true , Intentos : int64(camionp.regi[cargaEntregada[i]].intentos)})
		     		if err != nil { //en caso de que falle el reporte reintentamos
		     			fmt.Println("could not greet:", err)
		     			fmt.Println("Reintentando")
		     			i = i - 1
		     		}
		     	}
		     	for j := 0; j < len(cargaNoEntregada); j++ {//Reportamos los no entregados
		     		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			        defer cancel()
		     		_ , err := c.Report(ctx, &pb.ReportDelivery{IdPaquete : camionp.regi[cargaNoEntregada[j]].idPaquete,
		     		 Entregado : false ,Intentos : int64(camionp.regi[cargaNoEntregada[j]].intentos)})
		     		if err != nil { //en caso de que falle el reporte reintentamos
		     			fmt.Println("could not greet:", err)
		     			fmt.Println("Reintentando")
		     			j = j - 1
		     		}
		    	}
		    	cargaEntregada = []int{} //Vaciamos cargaEntregada y cargaNoEntregada
	            cargaNoEntregada = []int{}
			}
	    }
    }
}


 func main() {

 	waitFor2 := 30 

 	waitForDeliverPack := 10

 	fmt.Println(waitForDeliverPack)

	if len(os.Args) > 1 {
		wait,err := strconv.Atoi(os.Args[1])
		if err != nil {
			wait = 30
		}
		waitFor2 = wait
	}

	fmt.Println(waitFor2)

	wg.Add(1)
	go camionLaborando("normal",float64(waitFor2),float64(waitForDeliverPack)) 

	time.Sleep(6 * time.Second)

    wg.Add(1)
	go camionLaborando("retail",float64(waitFor2),float64(waitForDeliverPack))

	time.Sleep(6 * time.Second)

    wg.Add(1)
	go camionLaborando("retail",float64(waitFor2),float64(waitForDeliverPack))

    wg.Wait()
	fmt.Println(waitFor2)
}