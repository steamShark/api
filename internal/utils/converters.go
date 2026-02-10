package utils

/*
File to hold the functions that receive avalue and returns the pointer
*/

func ReturnPointerBool(boolean bool) *bool {
	return &boolean
}
