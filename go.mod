module plantbot

go 1.19

require tinygo.org/x/drivers v0.23.0

require (
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/text v0.3.6 // indirect
)

// bacause of wifinina Stop() function
replace tinygo.org/x/drivers => ../tinygo-drivers
