module github.com/ramalho/sets-in-go/setgm

// Generic methods require Go 1.27. Until 1.27.0 is released the toolchain
// line pins the release candidate, which the go command finds as the
// go1.27rc2 binary in PATH; without it, an older go tries to download a
// go1.27.0 that does not exist yet.
toolchain go1.27rc2

go 1.27

require github.com/ramalho/sets-in-go/vendored v0.0.0

replace github.com/ramalho/sets-in-go/vendored => ../vendored
