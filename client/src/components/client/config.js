const config = require('./config.json');
const host = config && config.APIHost ? config.APIHost : "http://localhost";
const port = config && config.APIPort ? (":" + config.APIPort) : "";
const baseUrl = host + port;
const useDefaultStyle = config && config.UseDefaultStyle ? config.UseDefaultStyle : true;
const moderation = config && config.Moderation ? config.Moderation : false;
const pageSize = config && config.PageSize ? config.PageSize : 5;
module.exports = {
    url: baseUrl,
    useDefaultStyle:useDefaultStyle,
    moderation:  moderation,
    pageSize: pageSize
}