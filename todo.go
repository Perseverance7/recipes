package todo

type TodoList struct {
    Id          int    `json:"id" db: "id"`
    Title       string `json:"title" db: "title" binding:"required"`
    Description string `json:"description" db: "description"`
}

type UsersList struct {
    Id     int `json:"id"`
    UserId int `json:"user_id"`
    ListId int `json:"list_id"`
}

type TodoItem struct {
    Id          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
}

type ListsItem struct {
    Id     int `json:"id"`
    ListId int `json:"list_id"`
    ItemId int `json:"item_id"`
}
