package main

import (
	directoryapp "DirectoryApp/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

func main()	{
	support := flag.String("sp", "", "Especifique el nombre del archivo local que subira como soporte")
	flag.Parse()

	client_conection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error al conectarse al localohst")
	}
	defer client_conection.Close()

	conection := directoryapp.NewWorkServiceClient(client_conection)

	supportfile, err := ioutil.ReadFile(path.Join(*support))
	if err != nil {
		log.Println("Error al subir el archivo", err)
	}

	supportname := strings.Split(*support, "\\")

	reqtask := &directoryapp.TaskRequest{
		Task: &directoryapp.Task{
			Id: "123",
			Description: "Cambiar cableado de electricidad",
			Profession: "Electricista",
			Isdone: false,
		},
	}

	reqclient := &directoryapp.CreateClientRequest{
		Client: &directoryapp.Client{
			Id: "123",
			Name: "Rachid Moyse",
			Phone: "3144351189",
			Email: "rachidmoyse1@hotmail.com",
			Password: "Skills39",
			Usertype: "Client",
			Tasks: []*directoryapp.Task{
				&directoryapp.Task{
					Id: reqtask.Task.Id,
					Description: reqtask.Task.Description,
					Profession: reqtask.Task.Profession,
					Isdone: reqtask.Task.Isdone,
				},
			},
		},
	}

	reqworker := &directoryapp.CreateWorkerRequest{
		Worker: &directoryapp.Worker{
			Id: "123",
			Name: "Alyeth Perez",
			Phone: "3014723047",
			Email: "alypg05@hotmail.com",
			Password: "Thor123",
			Profession: "Business Administrator",
			Usertype: "Worker",
			Supports: []string{fmt.Sprintf("%s", supportfile)},
			Supportsname: []string{fmt.Sprintf("%s", supportname[7])},
		},
	}
	
	reqapply := &directoryapp.ApplyRequest{
		Workername:   "Alyeth Perez",
		Existingtask: &directoryapp.Task{
			Id:          reqtask.Task.Id,
			Description: reqtask.Task.Description,
			Profession:  reqtask.Task.Profession,
			Requester:   reqtask.Task.GetRequester(),
			Isdone:      reqtask.Task.GetIsdone(),
		},
	}

	donereq := &directoryapp.UpdatetoDoneRequest{
		Task: &directoryapp.Task{
			Id:          reqtask.Task.Id,
			Description: reqtask.Task.Description,
			Profession:  reqtask.Task.Profession,
			Requester:   reqtask.Task.GetRequester(),
			Isdone:      reqtask.Task.GetIsdone(),
		},
	}

	resclnt, err := conection.CreateClient(context.Background(), reqclient)
	if err != nil {
		log.Printf("Error al crear el Client: %v", err)
	}
	fmt.Println(resclnt)

	time.Sleep(1000 * time.Millisecond)

	reswork, err := conection.CreateWorker(context.Background(), reqworker)
	if err != nil {
		log.Printf("Error al crear el Worker: %v", err)
	}
	fmt.Println(reswork)

	time.Sleep(1000 * time.Millisecond)

	restsk, err := conection.CreateWorkRequest(context.Background(), reqtask)
	if err != nil {
		log.Printf("Error al crear la tarea: %v", err)
	}
	fmt.Println(restsk)

	time.Sleep(1000 * time.Millisecond)

	restapply, err := conection.ApplytoWork(context.Background(), reqapply)
	if err != nil {
		log.Printf("Error al aplicar a la tarea: %v", err)
	}
	fmt.Println(restapply)

	time.Sleep(100000 * time.Millisecond)

	doneresp, err := conection.UpdateTasktoDone(context.Background(), donereq)
	if err != nil {
		log.Printf("Error al marcar como done la tarea: %v", err)
	}
	fmt.Println(doneresp)


}
