import { h, Component } from "preact";
import style from "./style";
const useStyle = true;


function getStyle(c) {
  return getStyle ? style[c] : c
}

function formatDate(d) {
  var dd = new Date(d)
  return dd.toISOString().slice(0, 19).replace("T", " ")
}

function findFormIndex(id, context) {
  return context.state.forms.map(x => x.id).indexOf(id)
}
const handleStateChange = (http, context) => {
  if (http.readyState != 4) {
    return
  }
  if (http.status == 200) {
    var parsedResponse = JSON.parse(http.responseText)
    if (parsedResponse.length > 0) {
      context.setState({ loaded: true, comments: parsedResponse, threadId: parsedResponse[0].ThreadId })
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
      forms: [{
        id: -1,
        visible: true,
        author: "",
        comment: "",
        email: null,
        canReply: false
      }]

    }
    this.fetchComments = this.fetchComments.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleAuthorChange = this.handleAuthorChange.bind(this);
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handleNewCommentSubmit = this.handleNewCommentSubmit.bind(this);
  }
  handleAuthorChange(id, value) {
    var form = findFormIndex(id, this)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { author: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }

  }
  handleEmailChange(id, value) {
    var form = findFormIndex(id, this)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { email: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }
  }
  handleBodyChange(id, value) {
    var form = findFormIndex(id, this)
    if (form >= 0) {
      var updatedForm = Object.assign({}, this.state.forms[form], { comment: value });
      var forms = this.state.forms;
      forms[form] = updatedForm;
      this.setState({ forms })
    }
  }
  handleNewCommentSubmit(id, replyTo) {
    if (typeof window == "undefined") { return }
    // validation
    var formIndex = findFormIndex(id, this)
    if (formIndex < 0) {
      return
    }
    var form = this.state.forms[formIndex];
    if (form.author == "") {
      // TODO show stuff
      return
    }
    if (form.comment == "") {
      // TODO show stuff
      return
    }

    var http = new XMLHttpRequest();
    var url = "http://localhost:7777/v1/comments";
    http.open("POST", url, true);
    var context = this;
    http.onreadystatechange = function () {
      if (http.readyState == 4) {
        if (http.status == 204) {
          // submit success, show the comment in the list below
          var cm = context.state.comments
          var maxId = 0;
          for (var i = 0; i < cm.length; i++) {
            if (cm[i].Id > maxId) {
              maxId = cm[i].Id
            }
          }
          cm.push({
            ThreadId: context.state.threadId,
            Id: ++maxId,
            Body: context.state.forms[formIndex].comment,
            Author: context.state.forms[formIndex].author,
            Confirmed: false,
            CreatedAt: new Date(),
            DeletedAt: null,
            ReplyTo: null
          })
          var updatedForm = Object.assign({}, context.state.forms[formIndex], { email: context.state.forms[formIndex].email, author: "", body: "" });
          var forms = context.state.forms;
          forms[form] = updatedForm;
          context.setState({ comments: cm, forms })
        } else {
          // error 
        }
      }
    }
    http.send(JSON.stringify({ Body: context.state.forms[formIndex].comment, Author: context.state.forms[formIndex].author, Path: window.location.pathname, ReplyTo:  null }))
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

  render(props) {
    var commentsFiltered = this.state.comments.filter(x => x.ReplyTo == null)
    var commentDiv = <div class={getStyle("mouthful_no_comments")}>No comments yet!</div>

    if (commentsFiltered.length != 0) {
      commentDiv = commentsFiltered.map(comment => {
        var replies = this.state.comments.filter(x => x.ReplyTo === comment.Id).map(x => {
          return <div class={getStyle("mouthful_comment_reply")} key={"___comment" + x.Id}>
            <div class={getStyle("mouthful_author")}>{x.Author}</div>
            <div class={getStyle("mouthful_date")}>{formatDate(x.CreatedAt)}</div>
            <div class={getStyle("mouthful_comment_body")}>{x.Body}</div>
          </div>
        });
        return <div class={getStyle("mouthful_comment")} key={"___comment" + comment.Id}>
          <div class={getStyle("mouthful_author")}>{comment.Author}</div>
          <div class={getStyle("mouthful_date")}>{formatDate(comment.CreatedAt)}</div>
          <div class={getStyle("mouthful_comment_body")}>{comment.Body}</div>
          <div class={getStyle("mouthful_comment_replies")}>
            {replies}
          </div>
        </div>;
      })
    }

    var form = (<div class={getStyle("mouthful_form")}>
      <input
        class={getStyle("mouthful_author_input")}
        type="text" name="author"
        placeholder="Name (required)"
        value={this.state.forms[findFormIndex(-1, this)].author}
        onChange={(e) => this.handleAuthorChange(-1, e.target.value)}>
      </input>
      <input
        style={"display: none;"}
        type="text" name="email"
        placeholder="Email (required)"
        value={this.state.forms[findFormIndex(-1, this)].email}
        onChange={(e) => this.handleEmailChange(-1, e.target.value)}>
      </input>
      <textarea
        class={getStyle("mouthful_comment_input")}
        rows="5"
        name="commentBody"
        placeholder="Type comment here..."
        value={this.state.forms[findFormIndex(-1, this)].comment}
        onChange={(e) => this.handleBodyChange(-1, e.target.value)}>
      </textarea>
      <input class={getStyle("mouthful_submit")} type="submit" value="Submit" onClick={(e) => { this.handleNewCommentSubmit(-1, null) }}></input>
    </div>)

    return (
      <div class={getStyle("mouthful_wrapper")}>
        {form}
        {commentDiv}
      </div>
    );
  }
}
