package protobuf

import (
	"io"

	pb "github.com/amosehiguese/go-chore/proto/v1"
	"google.golang.org/protobuf/proto"
)

func Load(r io.Reader) ([]*pb.Chore, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var chores pb.Chores
	return chores.Chores, proto.Unmarshal(b, &chores)
}

func Flush(w io.Writer, chores []*pb.Chore) error {
	b, err := proto.Marshal(&pb.Chores{Chores: chores})
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}