import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import config from "./config"
import Form from "./form"
import FormWrapper from "./formWrapper"
import Comment from "./comment"
const useStyle = config.useDefaultStyle;

function getStyle(c) {
  return useStyle ? style[c] : c
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
        return Object.assign({}, x, { RepliesToLoad: config.pageSize })
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
      error: false,
      comments: [],
      threadId: 0,
      author: cookies.get("mouthful_author") ? cookies.get("mouthful_author") : "",
      showComments: config.pageSize,
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
  }
  incrementReplyCount(commentId) {
    var copiedComments = this.state.comments;
    var found = copiedComments.map(x => x.Id).indexOf(commentId)
    if (found >= 0) {
      copiedComments[found].RepliesToLoad += config.pageSize
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
    var url = config.url + "/v1/comments";
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
              var leftOvers = (totalReplies % config.pageSize) > 0 ? 1 : 0;
              cm[found].RepliesToLoad =  totalReplies * config.pageSize + leftOvers * config.pageSize;
            }
          } else {
            var totalComments = cm.filter(x => x.ReplyTo == null).length + 1;
            var leftOvers = (totalComments % config.pageSize) > 0 ? 1 : 0;
            toShow = totalComments * config.pageSize + leftOvers * config.pageSize;
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
            RepliesToLoad: config.pageSize
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
          setTimeout(() => context.focus(config.commentRefPrefix + parsedResponse.id), 100)
        }
      }
    }
    http.send(JSON.stringify({
      Body: comment,
      Author: author,
      Path: window.location.pathname,
      ReplyTo: replyTo
    }))
  }
  isFormVisible(id) {
    filtered = this.state.forms.filter(x=>x.id == id)
    return filtered[0].visible;
  }
  componentDidMount() {
    if (!this.state.loaded) {
      this.fetchComments()
    }
  }

  fetchComments() {
    if (typeof window == "undefined") { return }
    var context = this;
    var http = new XMLHttpRequest();
    var url = config.url + "/v1/comments?uri=" + encodeURIComponent(window.location.pathname);
    http.open("GET", url, true);
    http.onreadystatechange = function () {
      handleStateChange(http, context)
    }
    http.send()
  }
  render(props) {
    if (this.state.error == true) {
      return <div class={getStyle("mouthful_wrapper")}><div class={getStyle("mouthful_error")}>The comments are temporarily unavailable</div></div>
    }
    var commentsFiltered = this.state.comments.filter(x => x.ReplyTo == null);
    var commentDiv = <div class={getStyle("mouthful_no_comments")}>No comments yet!</div>
    var loadMoreComments = null;
    if (commentsFiltered.length != 0) {
      if (this.state.showComments && this.state.showComments > 0) {
        if (commentsFiltered.length > this.state.showComments) {
          loadMoreComments = <input
            class={getStyle("mouthful_reply_button")}
            onClick={() => { this.setState({ showComments: this.state.showComments + config.pageSize }) }}
            type="Submit"
            value="Show more comments" >
          </input>
        }
        commentsFiltered = commentsFiltered.slice(0, this.state.showComments);
      }
      commentDiv = commentsFiltered.map(comment => {
        var cmntsToFilter = this.state.comments.filter(x => x.ReplyTo === comment.Id);
        var loadMoreReplies = null;
        if (comment.RepliesToLoad && comment.RepliesToLoad > 0) {
          if (cmntsToFilter.length > comment.RepliesToLoad) {
            loadMoreReplies = <input
              class={getStyle("mouthful_reply_button")}
              onClick={() => { this.incrementReplyCount(comment.Id) }}
              type="Submit"
              value="Show more replies" >
            </input>
          }
          cmntsToFilter = cmntsToFilter.splice(0, comment.RepliesToLoad)
        }

        var replies = cmntsToFilter.map(x => {
          var formIndex = this.findFormIndex(x.Id);
          return <div class={getStyle("mouthful_comment_reply")} key={"___comment" + x.Id} tabindex="-1" ref={c => {
            this.refMap.set(config.commentRefPrefix + x.Id, c)
          }}>
            <Comment comment={x}/>
            <FormWrapper comment={x} flipFormVisibility={this.flipFormVisiblity} visible={this.state.forms[this.findFormIndex(x.Id)].visible}  author={this.state.author}  replyTo={x.Id} submitForm={this.submitForm}/>
          </div>
        });
        var formIndex = this.findFormIndex(comment.Id);
        return <div class={getStyle("mouthful_comment")} key={"___comment" + comment.Id} tabindex="-1" ref={c => {
          this.refMap.set(config.commentRefPrefix + comment.Id, c)
        }}>
          <Comment comment={comment}/>
          <FormWrapper comment={comment} flipFormVisibility={this.flipFormVisiblity} visible={this.state.forms[this.findFormIndex(comment.Id)].visible}  author={this.state.author}  replyTo={comment.Id} submitForm={this.submitForm}/>
          <div class={getStyle("mouthful_comment_replies")}>
            {replies}
            {loadMoreReplies}
          </div>
        </div>;
      })
    }

    return (
      <div class={getStyle("mouthful_wrapper")}>
        <Form id={-1} visible={this.state.forms[this.findFormIndex(-1)].visible} author={this.state.author} comment={""} replyTo={null} submitForm={this.submitForm} />
        {commentDiv}
        {loadMoreComments}
      </div>
    );
  }
}
