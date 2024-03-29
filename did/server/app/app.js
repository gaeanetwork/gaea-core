const express = require('express')
const axios = require('axios')
var bodyParser = require('body-parser')
var session = require('express-session')

const app = express()
const port = 8081

app.use(express.static('public'), session({
	secret: 'cookie-secret-key',
	resave: false,
	saveUninitialized: true
}))
.use(function (req, res, next) {

	if (['/', '/login', '/verify'].indexOf(req.path) >= 0) {
		return next()
	}

	if (!req.session.did) {
		return res.status(401).send("Unauthorized, please log in.")
	}

	next()
})

app.all('*', function (req, res, next) {
	res.header('Access-Control-Allow-Origin', '*');
	res.header('Access-Control-Allow-Headers', 'Content-Type, Content-Length, Authorization, Accept, X-Requested-With');
	res.header('Access-Control-Allow-Methods', 'PUT, POST, GET, DELETE, OPTIONS');
	res.header("Content-Type", "text/html; charset=utf-8");

	if (req.method == 'OPTIONS') {
		res.send(200);
	}
	else {
		next();
	}
});

// serve index.html as the home page
app.get('/', function (req, res) {
	console.log(req.session.id)
	res.sendFile('index.html', { root: __dirname });
})

app.get('/login', function (req, res) {
	redirectUrl = 'http://localhost:8081/verify'
	// XXX TODO: sign message before send
	authRequest = {
		client_id: redirectUrl,
		nonce: req.session.id,
	  }
	  res.status(200).send(authRequest)
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
			res.header("Content-Type", "application/json;charset=utf-8");
			res.status(200).json({result:true});
		} else {
			res.header("Content-Type", "application/json;charset=utf-8");
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