# wip

**Build steps**

    go build main.go

**Steps to run the code**

    ./main <filePath>

**Steps to run unit tests**

    go test -v

**Sample Outputs:**

    $ # create build
    $ go build main.go

    $ # run the code
    $ ./main ./samplelogfiles/log1.txt
    ALICE99 4 240
    CHARLIE 3 37

    $ # run unit tests
    $ go test -v
    === RUN   TestParseLine
    --- PASS: TestParseLine (0.00s)
    === RUN   TestFetchLatestTime
    --- PASS: TestFetchLatestTime (0.00s)
    === RUN   TestFetchEarliestTime
    --- PASS: TestFetchEarliestTime (0.00s)
    === RUN   TestGetLines
    --- PASS: TestGetLines (0.00s)
    === RUN   TestGenerateReport
    --- PASS: TestGenerateReport (0.00s)
    PASS
    ok      wip     0.002s
