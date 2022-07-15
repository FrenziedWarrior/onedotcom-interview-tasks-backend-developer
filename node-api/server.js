const express = require("express");
const cors = require("cors");

const app = express();

var corsOptions = {
    methods: ['GET', 'POST', 'PATCH']
};

app.use(cors(corsOptions));

// parse requests of content-type - application/json
app.use(express.json());

// parse requests of content-type - application/x-www-form-urlencoded
app.use(express.urlencoded({ extended: true }));

const db = require("./app/models");
db.sequelize.sync()
.then(() => {
    console.log("Synced db.");
})
.catch((err) => {
    console.log("Failed to sync db: " + err.message);
});

// defining blacklist for IPs 
var blackListIPs = ['172.16.0.2'];

// geting client IP
var getClientIp = function(req) {
    var ipAddress = req.connection.remoteAddress;
    if (!ipAddress) {
        return '';
    }

    // convert from "::ffff:192.0.0.1"  to "192.0.0.1"
    if (ipAddress.substr(0, 7) == "::ffff:") {
        ipAddress = ipAddress.substr(7)
    }
    return ipAddress;
};

// Blocking Client IP, if it is in the blacklist
app.use(function(req, res, next) {
    var ipAddress = getClientIp(req);
    if (blackListIPs.indexOf(ipAddress) === -1) {
        next();
    } else {
        res.send({message: "Request is unauthorized!"})
    }
});

// To check request headers before proceeding to route requests
var middleware = {
    requestHeaderValidator: function(req, res, next) {
        var contentTypeHeader = req.headers['content-type'];
        if (contentTypeHeader == 'application/json') {
            next();
        } else {
            res.status(401).json({message: 'Request is malformed. Please check Content-Type header in request!'});
        }
    }
}

app.use(middleware.requestHeaderValidator)

require("./app/routes/plugin.routes")(app);

// set port, listen for requests
const PORT = process.env.PORT || 4545;

app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}.`);
});