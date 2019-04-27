.PHONY: all

all: test_example test_homework cleanup gendata

test_example:
	go test -v -run=TestExampleURLTop

test_homework:
	go test -v -run=TestURLTop

cleanup:
	go test -v -run=TestCleanData

gendata:
	go test -v -run=TestGenData

makeprofile:
	go test -v -bench=. -cpuprofile=cpuprofile -run=TestURLTop

makeprofile_example:
	go test -v -bench=. -cpuprofile=cpuprofile_example -run=TestExampleURLTop
