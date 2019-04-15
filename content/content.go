package content

type Content interface {
	Line(index int) (string, error)
	Count() int
}
