const config = require('./config.json');
const host = config && config.APIHost ? config.APIHost : "http://localhost";
const port = config && config.APIPort ? (":" + config.APIPort) : "";
const baseUrl = host + port;
const useDefaultStyle = config && config.UseDefaultStyle ? config.UseDefaultStyle : true;
const moderation = config && config.Moderation ? config.Moderation : false;
const pageSize = config && config.PageSize ? config.PageSize : 0;
const maxMessageLength = config && config.MaxCommentLength ? config.MaxCommentLength : 0;
const authorInputRefPrefix = "__mouthful_author_input_";
const commentInputRefPrefix = "__mouthful_comment_input_";
const commentRefPrefix = "__mouthful_comment_";
module.exports = {
    url: baseUrl,
    useDefaultStyle:useDefaultStyle,
    moderation:  moderation,
    pageSize: pageSize,
    maxMessageLength: maxMessageLength,
    authorInputRefPrefix,
    commentInputRefPrefix,
    commentRefPrefix,
}