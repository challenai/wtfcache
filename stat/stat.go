package stat

type Stat interface {
	Count() int64
	GetSz() int64
}
