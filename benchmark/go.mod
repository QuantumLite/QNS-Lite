module example.com/benchmark

go 1.15

replace example.com/profile => ../profile

replace example.com/graph => ../graph

replace example.com/config => ../config

replace example.com/path => ../path

replace example.com/quantum => ../quantum

replace example.com/request => ../request

replace example.com/log => ../log

require (
	example.com/config v0.0.0-00010101000000-000000000000
	example.com/graph v0.0.0-00010101000000-000000000000
	example.com/profile v0.0.0-00010101000000-000000000000
)
