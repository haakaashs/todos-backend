package handler

import (
	"context"
	"log"

	"connectrpc.com/connect"
	v1 "github.com/haakaashs/todos-backend/gen/protos/todos/v1"
	gen "github.com/haakaashs/todos-backend/gen/protos/todos/v1/todosv1connect"
	"github.com/haakaashs/todos-backend/internal/helper"
	"github.com/haakaashs/todos-backend/internal/model"
	"github.com/haakaashs/todos-backend/internal/service"
)

// TodosServiceHandler handles the TodoService gRPC requests.
type TodosServiceHandler struct {
	gen.UnimplementedTodosServiceHandler
	service *service.Service
}

// NewTodosServiceHandler creates a new TodosServiceHandler.
func NewTodosServiceHandler(service *service.Service) *TodosServiceHandler {
	return &TodosServiceHandler{service: service}
}

// Create implements the Create method of the TodoServiceHandler interface.
func (h *TodosServiceHandler) Create(ctx context.Context, req *connect.Request[v1.CreateRequest]) (*connect.Response[v1.CreateResponse], error) {
	log.Default().Println("Create todo method called")

	todo, err := h.service.Create(ctx, req.Msg.Title)
	if err != nil {
		return nil, err
	}

	res := &v1.Todo{}
	err = helper.TransformStruct(todo, res)
	if err != nil {
		return nil, err
	}

	log.Default().Println("Successfully created todo item")
	return connect.NewResponse(&v1.CreateResponse{Todo: res}), nil
}

// Get implements the Get method of the TodoServiceHandler interface.
func (h *TodosServiceHandler) Get(ctx context.Context, req *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error) {
	log.Default().Println("Get todo method called")

	todo, err := h.service.Get(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	res := &v1.Todo{}
	err = helper.TransformStruct(todo, res)
	if err != nil {
		return nil, err
	}

	log.Default().Println("Successfully fetched todo item")
	return connect.NewResponse(&v1.GetResponse{Todo: res}), nil
}

// Update implements the Update method of the TodoServiceHandler interface.
func (h *TodosServiceHandler) Update(ctx context.Context, req *connect.Request[v1.UpdateRequest]) (*connect.Response[v1.UpdateResponse], error) {
	log.Default().Println("Update todo method called")

	domainModel := &model.Todo{}
	err := helper.TransformStruct(req.Msg, domainModel)
	if err != nil {
		return nil, err
	}

	todo, err := h.service.Update(ctx, domainModel)
	if err != nil {
		return nil, err
	}

	res := &v1.Todo{}
	err = helper.TransformStruct(todo, res)
	if err != nil {
		return nil, err
	}

	log.Default().Println("Successfully updated todo item")
	return connect.NewResponse(&v1.UpdateResponse{Todo: res}), nil
}

// Delete implements the Delete method of the TodoServiceHandler interface.
func (h *TodosServiceHandler) Delete(ctx context.Context, req *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	log.Default().Println("Delete todo method called")

	err := h.service.Delete(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	log.Default().Println("Successfully deleted todo item")
	return connect.NewResponse(&v1.DeleteResponse{}), nil
}

// List implements the List method of the TodoServiceHandler interface.
func (h *TodosServiceHandler) List(ctx context.Context, req *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	log.Default().Println("List todos method called")

	todos, err := h.service.List(ctx)
	if err != nil {
		return nil, err
	}

	var resTodos []*v1.Todo
	err = helper.TransformStruct(todos, &resTodos)
	if err != nil {
		return nil, err
	}

	log.Default().Println("Successfully Listed todo items")
	return connect.NewResponse(&v1.ListResponse{Todos: resTodos}), nil
}
