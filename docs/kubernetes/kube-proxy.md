# kube-proxy 源码分析


### mode: iptables

```
// Table represents different iptable like filter,nat, mangle and raw
type Table string

const (
	// TableNAT represents the built-in nat table
	TableNAT Table = "nat"
	// TableFilter represents the built-in filter table
	TableFilter Table = "filter"
	// TableMangle represents the built-in mangle table
	TableMangle Table = "mangle"
)
