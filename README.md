# Signal Assessment

## Build

`docker build . -t assessment`

## Run

`docker run  -p 5000:5000 -e PORT=5000 -d assessment`

## Try it out

`./client.sh coinbasepro MATIC-BTC 1d`

## Points of improvement

- Split up the codebase into smaller chunks
- Expand the test suite to cover edge cases (got a bit lazy there :D)
- (maybe) return JSON, seems like overkill for an enum with 3 possible values though
- Parameterize the types of SMA(n) a client wants to compare