module example.com/request

go 1.15

replace example.com/graph => ../graph

replace example.com/config => ../config

require (
	example.com/config v0.0.0-00010101000000-000000000000
	example.com/graph v0.0.0-00010101000000-000000000000
)
