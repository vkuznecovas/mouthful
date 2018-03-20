const config = require('./config.json');
const host = config && config.APIHost ? config.APIHost : "http://localhost";
const port = config && config.APIPort ? (":" + config.APIPort) : "";
const baseUrl = host + port;
module.exports = {
    url: baseUrl
}