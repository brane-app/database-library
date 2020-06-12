package monkebase

type User interface{}

func WriteUser(user User) (err error) {
    return
}

func ReadSingleUser(ID string) (user User, exists bool, err error)  {
    return
}
