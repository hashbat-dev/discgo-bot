package todo

type ToDoType struct {
	Name        string
	Abbreviator string
	Description string
}

var ToDoTypes []ToDoType = []ToDoType{
	{
		Name:        "PRJFEM",
		Abbreviator: "P",
		Description: "Project",
	},
	{
		Name:        "SWKFEM",
		Abbreviator: "W",
		Description: "Small Works",
	},
	{
		Name:        "SUPFEM",
		Abbreviator: "U",
		Description: "Support/Bugs",
	},
}
