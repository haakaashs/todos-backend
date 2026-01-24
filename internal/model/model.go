package model

type Todo struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateRequest struct {
	Title string `json:"title"`
}

type CreateResponse struct {
	Todo *Todo `json:"todo"`
}

type GetRequest struct {
	Id string `json:"id"`
}

type GetResponse struct {
	Todo *Todo `json:"todo"`
}

type ListResponse struct {
	Todos []*Todo `json:"todos"`
}

type UpdateRequest struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UpdateResponse struct {
	Todo *Todo `json:"todo"`
}

type DeleteRequest struct {
	Id string `json:"id"`
}
