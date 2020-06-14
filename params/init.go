package params

import "flag"

var DeleteOnAnswer = flag.Bool("deleteOnAnswer", false, "set to delete command message as it is answered")
var DbFile = flag.String("db", "EDFCBank.db", "defines db file")

func init() {
	flag.Parse()
}
