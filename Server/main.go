package main

import (
	directoryapp "DirectoryApp/proto"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)


type server struct {
	directoryapp.UnimplementedWorkServiceServer
}
var collection *mongo.Collection
var client *mongo.Client

type User struct {
	ID string `bson:"id,required"`
	Name string `bson:"name,required"`
	Phone string `bson:"phone,required"`
	Email string `bson:"email,required"`
	Password string `bson:"password,required"`
	Usertype string `bson:"usertype,required"`
	Tasks Task `bson:"tasks,omitempty"`
}

type Worker struct {
	ID string `bson:"id,required"`
	Name string `bson:"name,required"`
	Phone string `bson:"phone,required"`
	Email string `bson:"email,required"`
	Tasks Task `bson:"tasks,omitempty"`
	Password string `bson:"password,required"`
	Supports byte `bson:"supports,omitempty"`
	Usertype string `bson:"usertype,required"`
}

type Task struct {
	ID string `bson:"id,required"`
	Description string `bson:"description,required"`
	Profession string `bson:"profession,required"`
	IsDone bool `bson:"is_done,required"`
}

func (s *server) CreateClient(ctx context.Context, req *directoryapp.CreateClientRequest) (*directoryapp.CreateClientResponse, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientid := uuid.NewV4().String()
	name := req.GetClient().GetName()
	phone := req.GetClient().GetPhone()
	email := req.GetClient().GetEmail()
	passwd := req.GetClient().GetPassword()
	tasks := req.GetClient().GetTasks()

	for _, task := range tasks {
		if nil == task {
			clnt := User{
				ID: clientid,
				Name: name,
				Phone: phone,
				Email: email,
				Password: passwd,
				Usertype: req.GetClient().GetUsertype(),
			}

			res, err := collection.InsertOne(ctx, clnt)
			if err != nil {
				log.Fatalln("Error al almacenar el usuario en la DB")
				return nil, err
			}

			resp := &directoryapp.CreateClientResponse{
				Client: fmt.Sprintf("Usuario insertado con ID: %v", res.InsertedID),
			}

			return resp, nil
		} else {
			clnt := User{
				ID: clientid,
				Name: name,
				Phone: phone,
				Email: email,
				Password: passwd,
				Usertype: req.GetClient().GetUsertype(),
				Tasks: Task{
					ID: uuid.NewV4().String(),
					Description: task.GetDescription(),
					Profession: task.GetProfession(),
					IsDone: task.GetIsdone(),
				},
			}

			res, err := collection.InsertOne(ctx, clnt)
			if err != nil {
				log.Fatalln("Error al almacenar el usuario en la DB")
				return nil, err
			}

			resp := &directoryapp.CreateClientResponse{
				Client: fmt.Sprintf("Usuario insertado con ID: %v", res.InsertedID),
			}
			return resp, nil
		}
	}
	return nil, nil
}

func (s *server) CreateWorker(ctx context.Context, req *directoryapp.CreateWorkerRequest) (*directoryapp.CreateWorkerResponse, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	workerid := uuid.NewV4().String()
	name := req.GetWorker().GetName()
	phone := req.GetWorker().GetPhone()
	email := req.GetWorker().GetEmail()
	passwd := req.GetWorker().GetPassword()
	supports := req.GetWorker().GetSupports()
	supportsname := req.GetWorker().GetSupportsname()

	for _, supportname := range supportsname {
		path, err :=  os.Getwd()
		file, err := os.Create(path+"\\"+supportname)
		if err != nil {
			log.Println("Error al crear el archivo del backup")
		}
		defer file.Close()

		for _, support := range supports {
			fmt.Sprintf("%v", support)
			nbyte, err := file.Write([]byte(support))
			if err != nil {
				log.Println("Error al escribir el archivo de backup")
			}

			log.Printf("Bytes scritos para el soporte: %v", nbyte)

		}
	}

	wrk := Worker{
		ID: workerid,
		Name: name,
		Phone: phone,
		Email: email,
		Password: passwd,
		Usertype: req.GetWorker().GetUsertype(),
	}

	res, err := collection.InsertOne(ctx, wrk)
	if err != nil {
		log.Fatalln("Error al almacenar el usuario en la DB")
		return nil, err
	}

	resp := &directoryapp.CreateWorkerResponse{
		Worker: fmt.Sprintf("Worker insertado con ID: %v", res.InsertedID),
	}

	return resp, nil
}

func (s *server) CreateWorkRequest(ctx context.Context, req *directoryapp.TaskRequest) (*directoryapp.TaskResponse, error) {
	task := req.GetTask()
	taskid := uuid.NewV4().String()

	filter := bson.D{{"name", "Rachid Moyse"}}
	var userCreator User
	err := collection.FindOne(context.TODO(), filter).Decode(&userCreator)
	if err != nil {
		log.Println("Error al encontrar el usuario filtrado")
	}

	tsk := Task{
		ID: userCreator.Tasks.ID,
		Description: userCreator.Tasks.Description,
		Profession: userCreator.Tasks.Profession,
		IsDone: userCreator.Tasks.IsDone,
	}

	res, err := collection.InsertOne(context.Background(), tsk)
	if err != nil {
		log.Println("Error al crear la tarea")
		return nil, err
	}
	fmt.Println(res)

	resp := &directoryapp.TaskResponse{
		Taskid: taskid,
		Description: task.GetDescription(),
	}

	return resp, nil
}

func (s *server) ApplytoWork(ctx context.Context, req *directoryapp.ApplyRequest) (*directoryapp.ApplyResponse, error) {
	existingtask := req.GetExistingtask()
	workername := req.GetWorkername()
	var workerUser Worker
	filter := bson.M{"name":workername}
	update := bson.M{
		"$set": bson.M{"tasks":existingtask},
	}

	err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&workerUser)
	if err != nil {
		log.Println("Error al buscar el worker filtrado", err)
	}
 	res := &directoryapp.ApplyResponse{Confirmation: fmt.Sprintf("Worker %v updated with tasks %v", workerUser.ID, existingtask.Id)}

	return res, nil

}

func (s *server) UpdateTasktoDone(ctx context.Context, req *directoryapp.UpdatetoDoneRequest) (*directoryapp.UpdatetoDoneResponse, error) {
	donerequetask := req.GetTask()
	filter := bson.M{"description":donerequetask.GetDescription()}
	var taskdone Task
	update := bson.M{
		"$set": bson.M{"is_done":true},
	}

	err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&taskdone)
	if err != nil {
		log.Println("Error al buscar la tarea filtrada", err)
	}

	res := &directoryapp.UpdatetoDoneResponse{Confirmation: fmt.Sprintf("Task %v is Complete and Done", taskdone.ID)}

	return res, nil
}

func main() {
	dburi := "mongodb+srv://dbadmin:dbadminpasswordRachid123*@directoryappcluster.trwxu.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Database configuration
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Println(err)
	}
	collection = client.Database("directoryapp").Collection("users")

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
	}
	fmt.Println("Successfully connected and pinged to db")

	collection = client.Database("directoryapp").Collection("Tasks")

	// ####################################################### GRPC ####################################################
	// GRPC Server
	lst, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Imposible escuchar por el puerto")
	}

	log.Println("Servidor levantado!!")

	grpc_server := grpc.NewServer()
	directoryapp.RegisterWorkServiceServer(grpc_server, &server{})

	if err := grpc_server.Serve(lst); err != nil {
		log.Fatalln("Error al levantar el servidor")
	}
}
