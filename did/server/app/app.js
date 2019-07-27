const express = require('express')

const app = express()
const port = 8081

// serve index.html as the home page
app.get('/', function (req, res) {
	res.sendFile('index.html', { root: __dirname })
})

// start server
app.listen(port, () => console.log(`Example app listening on port ${port}!`))