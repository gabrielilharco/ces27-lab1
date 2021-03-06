package main

import (
	"github.com/pauloaguiar/ces27-lab1/mapreduce"
	"hash/fnv"
	"unicode"
	"strings"
	"strconv"
)

// mapFunc is called for each array of bytes read from the splitted files. For wordcount
// it should convert it into an array and parses it into an array of KeyValue that have
// all the words in the input.
func mapFunc(input []byte) (result []mapreduce.KeyValue) {
	// 	Pay attention! We are getting an array of bytes.
	//
	// 	To decide if a character is a delimiter of a word, use the following check:
	//		!unicode.IsLetter(c) && !unicode.IsNumber(c)
	//
	//	Map should also make words lower cased:
	//		strings.ToLower(string)

	wordMap := make(map[string]int)
	wordKeyVal := make([]mapreduce.KeyValue, 0)

	inputString := strings.ToLower(string(input[:])) + " "
	currentWord := ""

	for _, inputChar := range inputString {
		if unicode.IsLetter(inputChar) || unicode.IsNumber(inputChar) {
			currentWord += string(inputChar)
		} else if len(currentWord) > 0 {
			value, exists := wordMap[currentWord]
			if !exists {
				value = 0
			}
			wordMap[currentWord] = value+1
			currentWord = ""
		}
	}

	for key, value := range wordMap {
		wordKeyVal = append(wordKeyVal, mapreduce.KeyValue{Key: key, Value: strconv.Itoa(value)})
	}

	return wordKeyVal
}

// reduceFunc is called for each merged array of KeyValue resulted from all map jobs.
// It should return a similar array that summarizes all similar keys in the input.
func reduceFunc(input []mapreduce.KeyValue) (result []mapreduce.KeyValue) {
	// 	Maybe it's easier if we have an auxiliary structure? Which one?
	//
	// 	You can check if a map have a key as following:
	// 		if _, ok := myMap[myKey]; !ok {
	//			// Don't have the key
	//		}
	//
	// 	Reduce will receive KeyValue pairs that have string values, you may need
	// 	convert those values to int before being able to use it in operations.
	//  	strconv.Atoi(string_number)

	wordMap := make(map[string]int)
	wordKeyVal := make([]mapreduce.KeyValue, 0)
	
	for _, keyVal := range input {
		value, exists := wordMap[keyVal.Key]
		if !exists {
			value = 0
		}
		currrentValue, err := strconv.Atoi(keyVal.Value)
		if err != nil {
			currrentValue = 1
		}
		wordMap[keyVal.Key] = value + currrentValue
		
	}
	
	for key, value := range wordMap {
		wordKeyVal = append(wordKeyVal, mapreduce.KeyValue{Key: key, Value: strconv.Itoa(value)})
	}
	
	return wordKeyVal
}

// shuffleFunc will shuffle map job results into different job tasks. It should assert that
// the related keys will be sent to the same job, thus it will hash the key (a word) and assert
// that the same hash always goes to the same reduce job.
// http://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
func shuffleFunc(task *mapreduce.Task, key string) (reduceJob int) {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() % uint32(task.NumReduceJobs))
}
