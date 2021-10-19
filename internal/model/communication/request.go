package communication

import "fmt"

type Request struct {
        Id uint64 `json:"-"`
        Service string `json:"-"`
        User string `json:"user"`
        Text string `json:"desc"`
}

func (r Request) String() string {
        return fmt.Sprintf("id: %d; user: %s; desc: %s",
                                r.Id, r.User, r.Text)
}
