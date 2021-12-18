# mLock

## Project Layout

Consider - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

## Test Dependencies

I'm not sure how to use `mockgen` properly with gomod, so I'm doing a `GO111MODULE=on go get github.com/golang/mock/mockgen@latest` and then adding it to the path `export PATH="$HOME/go/bin:$PATH"`. :(