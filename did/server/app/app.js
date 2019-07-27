const express = require('express')
const axios = require('axios')
var bodyParser = require('body-parser')

const app = express()
const port = 8081

app.use(express.static('public'))

app.all('*', function (req, res, next) {
	res.header('Access-Control-Allow-Origin', '*');
	res.header('Access-Control-Allow-Headers', 'Content-Type, Content-Length, Authorization, Accept, X-Requested-With');
	res.header('Access-Control-Allow-Methods', 'PUT, POST, GET, DELETE, OPTIONS');
	res.header("Content-Type", "application/json;charset=utf-8");

	if (req.method == 'OPTIONS') {
		res.send(200);
	}
	else {
		next();
	}
});

// serve index.html as the home page
app.get('/', function (req, res) {
	res.sendFile('index.html', { root: __dirname });
})

var textParser = bodyParser.text({ type: 'application/json' });
axios.defaults.baseURL = 'http://localhost:8082';
axios.defaults.headers.post['Content-Type'] = 'application/json';
app.post('/verify', textParser, function (req, res) {
	console.log('>>>>>>>>>>>>>>: <', req.body);
	axios.post('/verify', {sign: req.body['sign']})
	.then(function (response) {
		console.log('<<<<<<<<<<:', response.data);
		if (response.data.result == true) {
			res.status(200).json({result:true});
		} else {
			res.status(200).json({result:false});
		}
	})
	.catch(function (error) {
		console.log(error);
		res.status(401).send('Unable to verify authentication response: ' + error)
	})
})

// start server
app.listen(port, () => console.log(`Example app listening on port ${port}!`))