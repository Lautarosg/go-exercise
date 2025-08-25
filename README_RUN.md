# Running the Service

## Option 1: Inside Module Directory (Recommended)
```
cd go-exercise
go run .
```
Then POST to:
```
curl -X POST http://localhost:8080/api/v1/ltp -d '{"pairs":["BTC/USD","BTC/EUR"]}' -H 'Content-Type: application/json'
```

## Option 2: From Parent Directory Using Workspace
Create a `go.work` one level above (already handled if present):
```
go work init ./go-exercise
```
Then:
```
go run go-exercise
```

## Example Request Bodies
Empty list (returns empty ltp array):
```
{"pairs":[]}
```
Specific pairs:
```
{"pairs":["BTC/USD","BTC/CHF"]}
```
