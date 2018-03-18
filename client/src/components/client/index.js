import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
const useStyle = true;
const moderationEnabled = false;
const authorInputRefPrefix = "__mouthful_author_input_";
const commentInputRefPrefix = "__mouthful_comment_input_";
const commentRefPrefix = "__mouthful_comment_";
const defaultComments = 5;
function getStyle(c) {
  return useStyle ? style[c] : c
}

function formatDate(d) {
  var dd = new Date(d)
  return timeago(dd)
  return dd.toISOString().slice(0, 19).replace("T", " ")
}


const handleStateChange = (http, context) => {
  if (http.readyState != 4) {
    return
  }
  if (http.status == 200) {
    var parsedResponse = JSON.parse(http.responseText)

    if (parsedResponse.length > 0) {
      var forms = context.state.forms;
      var authorCookieValue = cookies.get("mouthful_author")
      var formsToAppend = parsedResponse.map(x => {
        return {
          id: x.Id,
          visible: false,
          author: authorCookieValue ? authorCookieValue : "",
          comment: "",
          email: null,
          replyTo: x.ReplyTo ? x.ReplyTo : x.Id,
        }
      })
      forms = forms.concat(formsToAppend)
      parsedResponse = parsedResponse.map(x => {
        return Object.assign({}, x, { RepliesToLoad: defaultComments })
      })
      context.setState({ loaded: true, comments: parsedResponse, threadId: parsedResponse[0].ThreadId, forms })
    } else {
      context.setState({ loaded: true, comments: [] })
    }
  } else if (http.status == 404) {
    context.setState({ loaded: true, comments: [] })
  } else {
    context.setState({ loaded: true })
    console.log("error while fetching");
  }
}
export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: false,
      comments: [],
      threadId: 0,
      showComments: defaultComments,
      forms: [{
        id: -1,
        visible: true,
        author: "",
        comment: "",
        email: null,
        replyTo: null
      }],
    }
    this.fetchComments = this.fetchComments.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleAuthorChange = this.handleAuthorChange.bind(this);
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handleNewCommentSubmit = this.handleNewCommentSubmit.bind(this);
    this.getForm = this.getForm.bind(this);
    this.flipFormVisiblity = this.flipFormVisiblity.bind(this);
    this.findFormIndex = this.findFormIndex.bind(this);
    this.refMap = new Map();
    this.focus = this.focus.bind(this);
    this.incrementReplyCount = this.incrementReplyCount.bind(this);
  }
  incrementReplyCount(commentId) {
    var copiedComments = this.state.comments;
    var found = copiedComments.map(x => x.Id).indexOf(commentId)
    console.log(found);
    if (found >= 0) {
      copiedComments[found].RepliesToLoad += defaultComments
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
  flipFormVisiblity(index) {
    var forms = this.state.forms;
    forms[index].visible = !forms[index].visible
    this.setState({ forms })
  }
  handleAuthorChange(id, value) {
    var form = this.findFormIndex(id)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { author: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }

  }
  handleEmailChange(id, value) {
    var form = this.findFormIndex(id)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { email: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }
  }
  handleBodyChange(id, value) {
    var form = this.findFormIndex(id)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { comment: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }
  }
  handleNewCommentSubmit(id) {
    if (typeof window == "undefined") { return }
    // validation
    var formIndex = this.findFormIndex(id)
    if (formIndex < 0) {
      return
    }
    var form = this.state.forms[formIndex];
    var authorCopy = form.author.replace(/\s/g,'');

    if (authorCopy == "" || authorCopy.length < 3) {
      this.focus(authorInputRefPrefix + form.id)
      return
    }

    var commentCopy = form.comment.replace(/\s/g,'');
    if (commentCopy == "" || commentCopy.length < 3) {
      this.focus(commentInputRefPrefix + form.id)
      return
    }

    var http = new XMLHttpRequest();
    var url = "http://localhost:7777/v1/comments";
    http.open("POST", url, true);
    var context = this;
    http.onreadystatechange = function () {
      if (http.readyState == 4) {
        if (http.status == 200) {
          // submit success, show the comment in the list below
          var cm = context.state.comments;
          var maxId = 0;
          for (var i = 0; i < cm.length; i++) {
            if (cm[i].Id > maxId) {
              maxId = cm[i].Id
            }
          }
          var toShow = context.state.showComments;
          if (context.state.forms[formIndex].replyTo != null) {
            var found = cm.map(x => x.Id).indexOf(context.state.forms[formIndex].replyTo)
            if (found > -1){
              var totalReplies = cm.filter(x => x.ReplyTo == context.state.forms[formIndex].replyTo).length + 1;
              var leftOvers = (totalReplies % defaultComments) > 0 ? 1 : 0;
              cm[found].RepliesToLoad =  totalReplies * defaultComments + leftOvers * defaultComments;
            }
          } else {
            var totalComments = cm.filter(x => x.ReplyTo == null).length + 1;
            var leftOvers = (totalComments % defaultComments) > 0 ? 1 : 0;
            toShow = totalComments * defaultComments + leftOvers * defaultComments;
          }
          var parsedResponse = JSON.parse(http.responseText)
          cm.push({
            ThreadId: context.state.threadId,
            Id: ++maxId,
            Body: parsedResponse.body,
            Author: context.state.forms[formIndex].author,
            Confirmed: false,
            CreatedAt: new Date(),
            DeletedAt: null,
            ReplyTo: context.state.forms[formIndex].replyTo,
            RepliesToLoad: defaultComments
          })
          var updatedForm = Object.assign(
            {},
            context.state.forms[formIndex],
            { email: context.state.forms[formIndex].email, author:  context.state.forms[formIndex].author, comment: "", visible: id == -1 ? true : false }
          );
          var authorCookieValue = cookies.get("mouthful_author")
          if (!authorCookieValue) {
            cookies.set("mouthful_author", context.state.forms[formIndex].author, 365)
          }
          var forms = context.state.forms;
          forms[formIndex] = updatedForm;
          forms = forms.map(x => {
            return Object.assign({}, x, {author: context.state.forms[formIndex].author})
          })
          context.setState({ comments: cm, forms, showComments: toShow })
          setTimeout(() => context.focus(commentRefPrefix + maxId), 100)
        }
      }
    }
    http.send(JSON.stringify({
      Body: context.state.forms[formIndex].comment,
      Author: context.state.forms[formIndex].author,
      Path: window.location.pathname,
      ReplyTo: context.state.forms[formIndex].replyTo
    }))
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
    var url = "http://localhost:7777/v1/comments?uri=" + encodeURIComponent(window.location.pathname);
    http.open("GET", url, true);
    http.onreadystatechange = function () {
      handleStateChange(http, context)
    }
    http.send()
  }
  getForm(id) {
    var form = this.state.forms[this.findFormIndex(id)]
    if (!form) {
      return null
    }
    return (<div class={getStyle(form.visible ? "mouthful_form" : "mouthful_form_invisible")}>
      <input
        class={getStyle("mouthful_author_input")}
        type="text" name="author"
        placeholder="Name (required)"
        value={this.state.forms[this.findFormIndex(id)].author}
        ref={c => {
          this.refMap.set(authorInputRefPrefix + id, c)
        }}
        onChange={(e) => this.handleAuthorChange(id, e.target.value)}>
      </input>
      <input
        style={"display: none;"}
        type="text" name="email"
        placeholder="Email (required)"
        value={this.state.forms[this.findFormIndex(id)].email}
        onChange={(e) => this.handleEmailChange(id, e.target.value)}>
      </input>
      <textarea
        class={getStyle("mouthful_comment_input")}
        rows="3"
        name="commentBody"
        placeholder="Type comment here..."
        ref={c => {
          this.refMap.set(commentInputRefPrefix + id, c)
        }}
        value={this.state.forms[this.findFormIndex(id)].comment}
        onChange={(e) => this.handleBodyChange(id, e.target.value)}>
      </textarea>
      <input
        class={getStyle("mouthful_submit")}
        type="submit"
        value="Submit"
        onClick={(e) => { this.handleNewCommentSubmit(id) }}>
      </input>
    </div>)
  }
  render(props) {
    var commentsFiltered = this.state.comments.filter(x => x.ReplyTo == null);
    var commentDiv = <div class={getStyle("mouthful_no_comments")}>No comments yet!</div>
    var loadMoreComments = null;
    if (commentsFiltered.length != 0) {
      if (this.state.showComments && this.state.showComments > 0) {
        if (commentsFiltered.length > this.state.showComments) {
          loadMoreComments = <input
            class={getStyle("mouthful_reply_button")}
            onClick={() => { this.setState({ showComments: this.state.showComments + defaultComments }) }}
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
            this.refMap.set(commentRefPrefix + x.Id, c)
          }}>
            <div class={getStyle("mouthful_author")}>{x.Author}
              <span class={getStyle("mouthful_date")}>{formatDate(x.CreatedAt)}</span>
              {formIndex < 0 && moderationEnabled ? <span class={getStyle("mouthful_moderation")}>In queue for moderation</span> : null}
            </div>

            <div class={getStyle("mouthful_comment_body")} dangerouslySetInnerHTML={{ __html: x.Body }} />
            {
              formIndex < 0
                ? null
                : <input
                  class={getStyle("mouthful_reply_button")}
                  onClick={() => this.flipFormVisiblity(formIndex)}
                  type="submit"
                  value={this.state.forms[formIndex].visible ? "Close" : "Reply"}>
                </input>
            }
            {this.getForm(x.Id)}
          </div>
        });
        var formIndex = this.findFormIndex(comment.Id);
        return <div class={getStyle("mouthful_comment")} key={"___comment" + comment.Id} tabindex="-1" ref={c => {
          this.refMap.set(commentRefPrefix + comment.Id, c)
        }}>
          <div class={getStyle("mouthful_author")}>{comment.Author}
            <span class={getStyle("mouthful_date")}>{formatDate(comment.CreatedAt)}</span>
            {formIndex < 0 && moderationEnabled ? <span class={getStyle("mouthful_moderation")}>In queue for moderation</span> : null}
          </div>

          <div class={getStyle("mouthful_comment_body")} dangerouslySetInnerHTML={{ __html: comment.Body }} />
          {
            formIndex < 0
              ? null
              : <input
                class={getStyle("mouthful_reply_button")}
                onClick={() => this.flipFormVisiblity(formIndex)}
                type="submit"
                value={this.state.forms[formIndex].visible ? "Close" : "Reply"}>
              </input>
          }
          {this.getForm(comment.Id)}
          <div class={getStyle("mouthful_comment_replies")}>
            {replies}
            {loadMoreReplies}
          </div>
        </div>;
      })
    }

    var form = this.getForm(-1)

    return (
      <div class={getStyle("mouthful_wrapper")}>
        {form}
        {commentDiv}
        {loadMoreComments}
      </div>
    );
  }
}
