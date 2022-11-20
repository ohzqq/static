package idxgen

type Index struct {
	Body     string
	Meta     string
	Children []string `toml:"children"`
	Path     string   `toml:"path"`
}

type Meta struct {
	Title    string   `toml:"title"`
	Template string   `toml:"template"`
	Tags     []string `toml:"tags"`
	Cmd      `toml:"cmd"`
}

type Cmd struct {
	Bin  string   `toml:"bin"`
	Args []string `toml:"args"`
}
