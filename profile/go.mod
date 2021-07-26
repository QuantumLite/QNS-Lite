module example.com/profile

go 1.15

replace example.com/graph => ../graph

replace example.com/config => ../config

replace example.com/request => ../request

replace example.com/path => ../path

replace example.com/quantum => ../quantum

replace example.com/log => ../log

require (
	example.com/config v0.0.0-00010101000000-000000000000
	example.com/graph v0.0.0-00010101000000-000000000000
	example.com/path v0.0.0-00010101000000-000000000000
	example.com/quantum v0.0.0-00010101000000-000000000000
	example.com/request v0.0.0-00010101000000-000000000000
	example.com/log v0.0.0-00010101000000-000000000000
)
