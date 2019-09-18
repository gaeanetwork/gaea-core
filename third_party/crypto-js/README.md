# Unit Testing - qunit

## qunit install

`npm install -g qunit`

## qunit testing

```js
// QUnit defaults to looking for test files in `test`
qunit
// or you can also put them anywhere and then specify file paths or glob expressions
qunit 'tests/*-test.js'

// to see more details
qunit --help
```