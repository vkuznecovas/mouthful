import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import Form from "./form"
import FormWrapper from "./formWrapper"
import Comment from "./comment"
function sortComments(a, b) {
  return new Date(a.CreatedAt) - new Date(b.CreatedAt)
}
const handleStateChange = (http, context) => {
  if (http.readyState != 4) {
    return
  }
  if (http.status == 200) {
    var parsedResponse = JSON.parse(http.responseText)

    if (parsedResponse.length > 0) {
      var forms = context.state.forms;
      var formsToAppend = parsedResponse.map(x => {
        return {
          id: x.Id,
          visible: false
        }
      })
      forms = forms.concat(formsToAppend)
      parsedResponse = parsedResponse.map(x => {
        return Object.assign({}, x, { RepliesToLoad: context.state.config.pageSize })
      })
      context.setState({ loaded: true, comments: parsedResponse, threadId: parsedResponse[0].ThreadId, forms })
    } else {
      context.setState({ loaded: true, comments: [] })
    }
  } else if (http.status == 404) {
    context.setState({ loaded: true, comments: [] })
  } else {
    context.setState({ loaded: true, error: true })
    console.log("error while fetching");
  }
}
export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: false,
      configLoaded: false,
      error: false,
      hostUrl: "",
      config: {
        useDefaultStyle: false,
        moderation:  false,
        pageSize: 0,
        maxMessageLength: 0,
        maxAuthorLength: 0,
        authorInputRefPrefix: "__mouthful_author_input_",
        commentInputRefPrefix: "__mouthful_comment_input_",
        commentRefPrefix: "__mouthful_comment_",
      },
      comments: [],
      threadId: 0,
      author: cookies.get("mouthful_author") ? cookies.get("mouthful_author") : "",
      showComments: 0,
      forms: [{
        id: -1,
        visible: true,
      }],
    }
    this.fetchComments = this.fetchComments.bind(this);
    this.flipFormVisiblity = this.flipFormVisiblity.bind(this);
    this.findFormIndex = this.findFormIndex.bind(this);
    this.refMap = new Map();
    this.focus = this.focus.bind(this);
    this.incrementReplyCount = this.incrementReplyCount.bind(this);
    this.submitForm = this.submitForm.bind(this);
    this.isFormVisible = this.isFormVisible.bind(this);
    this.fetchConfig = this.fetchConfig.bind(this);
    this.getStyle = this.getStyle.bind(this);
  }
  getStyle(c) {
    return this.state.config.useDefaultStyle ? style[c] : c
  }
  incrementReplyCount(commentId) {
    var copiedComments = this.state.comments;
    var found = copiedComments.map(x => x.Id).indexOf(commentId)
    if (found >= 0) {
      copiedComments[found].RepliesToLoad += this.props.config.pageSize
      this.setState({ comments: copiedComments })
    }
  }
  focus(focusThis) {
    var tf = this.refMap.get(focusThis);
    if (tf) {
      tf.focus();
    }
  }
  findFormIndex(id) {
    return this.state.forms.map(x => x.id).indexOf(id)
  }
  flipFormVisiblity(id) {
    var forms = this.state.forms;
    var index = this.findFormIndex(id)
    forms[index].visible = !forms[index].visible
    this.setState({ forms })
  }
  submitForm(id, comment, author, email, replyTo) {
    var http = new XMLHttpRequest();
    var url = this.state.hostUrl + "/v1/comments";
    http.open("POST", url, true);
    var context = this;
    http.onreadystatechange = function () {
      if (http.readyState == 4) {
        if (http.status == 200) {
          // submit success, show the comment in the list below
          var cm = context.state.comments;
          var toShow = context.state.showComments;
          if (replyTo != null) {
            var found = cm.map(x => x.Id).indexOf(replyTo)
            if (found > -1){
              var totalReplies = cm.filter(x => x.ReplyTo == replyTo).length + 1;
              var leftOvers = (totalReplies % context.state.config.pageSize) > 0 ? 1 : 0;
              cm[found].RepliesToLoad =  totalReplies * context.state.config.pageSize + leftOvers * context.state.config.pageSize;
            }
          } else {
            var totalComments = cm.filter(x => x.ReplyTo == null).length + 1;
            var leftOvers = (totalComments % context.state.config.pageSize) > 0 ? 1 : 0;
            toShow = totalComments * context.state.config.pageSize + leftOvers * context.state.config.pageSize;
          }
          var parsedResponse = JSON.parse(http.responseText)
          
          cm.push({
            ThreadId: context.state.threadId,
            Id: parsedResponse.id,
            Body: parsedResponse.body,
            Author: author,
            Confirmed: false,
            CreatedAt: new Date(),
            DeletedAt: null,
            ReplyTo: replyTo,
            RepliesToLoad: context.state.config.pageSize
          })
          var forms = context.state.forms
          if (replyTo) {
            var index = context.findFormIndex(replyTo)
            forms[index].visible = false
          }
          forms = forms.concat([{
            id: parsedResponse.id,
            visible: false,
          }])
          
          var authorCookieValue = cookies.get("mouthful_author")
          if (!authorCookieValue) {
            cookies.set("mouthful_author", author, 365)
          }
          context.setState({ comments: cm, showComments: toShow, author: author, forms })
          setTimeout(() => context.focus(context.state.config.commentRefPrefix + parsedResponse.id), 100)
        }
      }
    }
    var path = window.location.pathname;
    if (this.state.pathPrefix) {
      path = this.state.pathPrefix + window.location.pathname;
    }
    var bod = {
      Body: comment,
      Author: author,
      Path: path,
      ReplyTo: replyTo
    }
    if (email != null) {
      bod.Email = email
    }
    http.send(JSON.stringify(bod))
  }
  isFormVisible(id) {
    filtered = this.state.forms.filter(x=>x.id == id)
    return filtered[0].visible;
  }
  componentDidMount() {
    var prefix = document.querySelector("#mouthful-comments").dataset.domain
    // remove the trailing slash if it's there
    if (prefix && prefix.endsWith("/")){
      prefix = prefix.substring(0, str.length-1);
    }
    this.setState({
      hostUrl: document.querySelector("#mouthful-comments").dataset.url,
      pathPrefix: prefix,
    })
    if (!this.state.configLoaded && this.state.hostUrl != "") {
      this.fetchConfig()
    }    
    if (!this.state.loaded && this.state.hostUrl != "") {
      this.fetchComments()
    }
  }
  fetchConfig() {
    if (typeof window == "undefined") { return }
    var context = this;
    var http = new XMLHttpRequest();
    var url = this.state.hostUrl + "/v1/client/config";
    http.open("GET", url, true);
    http.onreadystatechange = function () {
      if (http.readyState != 4) {
        return
      }
      if (http.status == 200) {
        var parsedResponse = JSON.parse(http.responseText)
        context.setState({ configLoaded: true, config: Object.assign(context.state.config, parsedResponse), showComments: parsedResponse.pageSize })
      } else {
        context.setState({ configLoaded: true, error: true })
        console.log("error while fetching config");
      }
    }
    http.send()
  }
  fetchComments() {
    if (typeof window == "undefined") { return }
    var context = this;
    var http = new XMLHttpRequest();
    var path = window.location.pathname;
    if (this.state.pathPrefix) {
      path = this.state.pathPrefix + window.location.pathname;
    }
    var url = this.state.hostUrl + "/v1/comments?uri=" + encodeURIComponent(path);
    http.open("GET", url, true);
    http.onreadystatechange = function () {
      handleStateChange(http, context)
    }
    http.send()
  }
  render(props) {
    if (this.state.error == true) {
      return <div class={this.getStyle("mouthful_wrapper")}><div class={this.getStyle("mouthful_error")}>The comments are temporarily unavailable</div></div>
    }
    if (!this.state.configLoaded || !this.state.loaded) {
      return <div class={this.getStyle("mouthful_wrapper")}><div class={this.getStyle("mouthful_no_comments")}>Loading...</div></div>
    }
    var commentsFiltered = this.state.comments.filter(x => x.ReplyTo == null).sort(sortComments);
    var commentDiv = <div class={this.getStyle("mouthful_no_comments")}>No comments yet!</div>
    var loadMoreComments = null;
    if (commentsFiltered.length != 0) {
      if (this.state.showComments && this.state.showComments > 0) {
        if (commentsFiltered.length > this.state.showComments) {
          loadMoreComments = <input
            class={this.getStyle("mouthful_reply_button")}
            onClick={() => { this.setState({ showComments: this.state.showComments + this.state.config.pageSize }) }}
            type="Submit"
            value="Show more comments" >
          </input>
        }
        commentsFiltered = commentsFiltered.slice(0, this.state.showComments);
      }
      commentDiv = commentsFiltered.map(comment => {
        var cmntsToFilter = this.state.comments.filter(x => x.ReplyTo === comment.Id).sort(sortComments);
        var loadMoreReplies = null;
        if (comment.RepliesToLoad && comment.RepliesToLoad > 0) {
          if (cmntsToFilter.length > comment.RepliesToLoad) {
            loadMoreReplies = <input
              class={this.getStyle("mouthful_reply_button")}
              onClick={() => { this.incrementReplyCount(comment.Id) }}
              type="Submit"
              value="Show more replies" >
            </input>
          }
          cmntsToFilter = cmntsToFilter.splice(0, comment.RepliesToLoad)
        }

        var replies = cmntsToFilter.map(x => {
          var formIndex = this.findFormIndex(x.Id);
          return <div class={this.getStyle("mouthful_comment_reply")} key={"___comment" + x.Id} tabindex="-1" ref={c => {
            this.refMap.set(this.state.config.commentRefPrefix + x.Id, c)
          }}>
            <Comment comment={x} config={this.state.config}/>
            <FormWrapper comment={x} config={this.state.config} flipFormVisibility={this.flipFormVisiblity} visible={this.state.forms[this.findFormIndex(x.Id)].visible}  author={this.state.author}  replyTo={comment.Id} submitForm={this.submitForm}/>
          </div>
        });
        var formIndex = this.findFormIndex(comment.Id);
        return <div class={this.getStyle("mouthful_comment")} key={"___comment" + comment.Id} tabindex="-1" ref={c => {
          this.refMap.set(this.state.config.commentRefPrefix + comment.Id, c)
        }}>
          <Comment comment={comment} config={this.state.config}/>
          <FormWrapper comment={comment} config={this.state.config} flipFormVisibility={this.flipFormVisiblity} visible={this.state.forms[this.findFormIndex(comment.Id)].visible}  author={this.state.author}  replyTo={comment.Id} submitForm={this.submitForm}/>
          <div>
            {replies}
            {loadMoreReplies}
          </div>
        </div>;
      })
    }

    return (
      <div class={this.getStyle("mouthful_wrapper")}>
        <Form id={-1} config={this.state.config} visible={this.state.forms[this.findFormIndex(-1)].visible} author={this.state.author} comment={""} replyTo={null} submitForm={this.submitForm} />
        {commentDiv}
        {loadMoreComments}
      </div>
    );
  }
}
