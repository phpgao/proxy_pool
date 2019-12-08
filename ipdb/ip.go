package ipdb

var (
	Db *City
)

func init() {
	r, e := newReaderFromGo("ipipfree.ipdb", &CityInfo{})
	if e != nil {
		panic(e)
	}
	Db = &City{
		reader: r,
	}
}
