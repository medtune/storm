package stormtf

type Query struct {
	Q []string
}

type QueryOption struct {
	MaxThreads   int8
	ProtoBufType string
	OutputType   string
	OutputName   string
	ResizeImage  bool
	ResizeDim    []int8
	VerboseLevel int8
	Queries      []Query
}

type StormTF struct {
	crawler      struct{}
	googleSearch struct{}
	writer       struct{}
	maxThread    int8
}

/*
stormtf --labels=cat;kity/dog --proto-format=features --resize=64*64 -o catdogs.tfrecord

stormtf --query=cat:cat+kitten/dog:dog --image=true
stormtf [-q --query] [-i --image] [-r --image-resize] [-o --output] [-v --verbose] [-p --proto-format]

*/

func New() *StormTF {
	return &StormTF{
		crawler:   struct{}{},
		maxThread: 8,
	}
}
