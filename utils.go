package main

func findIndex[T any](slice []T, matchFunc func(T) bool) int {
	for index, element := range slice {
		if matchFunc(element) {
			return index
		}
	}

	return -1 // not found
}

func removeElementByRange[T any](slice []T, from, to int) []T {
	return append(slice[:from], slice[to:]...)
}
