import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"

export default class Form extends Component {
  
  constructor(props) {
    super(props);
    this.state = {
      author: this.props.author,
      comment: this.props.comment,
      email: null,
      replyTo: this.props.replyTo ? this.props.replyTo : null
    }
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleAuthorChange = this.handleAuthorChange.bind(this);
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handleNewCommentSubmit = this.handleNewCommentSubmit.bind(this);
    this.refMap = new Map();
    this.focus = this.focus.bind(this);
    this.getStyle = this.getStyle.bind(this);
  }
  getStyle(c) {
    return this.props.config.useDefaultStyle ? style[c] : c
  }
  focus(focusThis) {
    var tf = this.refMap.get(focusThis);
    if (tf) {
      tf.focus();
    }
  }
  handleAuthorChange(value) {
    if (this.props.config.maxAuthorLength > 0) {
      var currentAuthor = this.state.author
      if (currentAuthor.length > (this.props.config.maxAuthorLength)) {
        if (value.length > currentAuthor.length) {
          // don't allow for extra characters, reset state to previous
          this.setState({ author: currentAuthor })
          return
        }
      }
    }
    this.setState({ author: value })
  }
  handleEmailChange(value) {
      this.setState({ email: value })
  }
  handleBodyChange(value) {
      if (this.props.config.maxMessageLength > 0) {
        var currentComment = this.state.comment
        if (currentComment.length > (this.props.config.maxMessageLength + 99)) {
          if (value.length > currentComment.length) {
            // don't allow for extra characters, reset state to previous
            this.setState({ comment: currentComment })
            return
          }
        }
      }
      this.setState({ comment: value })
  }

  handleNewCommentSubmit() {
    var authorCopy = this.state.author.replace(/\s/g,'');

    if (authorCopy == "" || authorCopy.length < 3) {
      this.focus(this.props.config.authorInputRefPrefix + this.props.id)
      return
    }

    var commentCopy = this.state.comment.replace(/\s/g,'');
    if (commentCopy == "" || commentCopy.length < 3) {
      this.focus(this.props.config.commentInputRefPrefix + this.props.id)
      return
    }

    if (this.props.config.maxMessageLength > 0) {
      if (this.state.comment.length > config.maxMessageLength) {
        this.focus(this.props.config.commentInputRefPrefix + this.props.id)
        return
      } 
    }
    this.props.submitForm(this.props.id, this.state.comment, this.state.author, this.state.email, this.state.replyTo)
    this.setState({comment: ""})
  }

  
  render(props) {
    var diff = this.props.config.maxMessageLength - this.state.comment.length;
    return (<div class={this.getStyle(this.props.visible ? "mouthful_form" : "mouthful_form_invisible")}>
      <input
        class={this.getStyle("mouthful_author_input")}
        type="text" name="author"
        placeholder="Name (required)"
        value={this.state.author}
        ref={c => {
          this.refMap.set(this.props.config.authorInputRefPrefix + this.props.id, c)
        }}
        onChange={(e) => this.handleAuthorChange(e.target.value)}
        onKeyUp={(e) => this.handleAuthorChange(e.target.value)}>
             
      </input>
      
      <input
        style={"display: none;"}
        type="text" name="email"
        placeholder="Email (required)"
        value={this.state.email}
        onChange={(e) => this.handleEmailChange(e.target.value)}>
      </input>
      <textarea
        class={this.getStyle("mouthful_comment_input")}
        rows="3"
        name="commentBody"
        placeholder="Type comment here..."
        ref={c => {
          this.refMap.set(this.props.config.commentInputRefPrefix + this.props.id, c)
        }}
        value={this.state.comment}
        onKeyUp={(e) => this.handleBodyChange(e.target.value)}
        onChange={(e) => this.handleBodyChange(e.target.value)}>
      </textarea>
     <div>
      <input
        class={this.getStyle("mouthful_submit")}
        type="submit"
        value="Submit"
        onClick={(e) => {this.handleNewCommentSubmit()}}>
      </input>
      {this.props.config.maxMessageLength > 0 ? <span class={diff > 0 ? this.getStyle("mouthful_word_counter") : this.getStyle("mouthful_word_counter_error")}>
                          {diff > 0 ? diff : diff * -1} {diff > 0 ? "characters left" : "characters too many"}
                      </span> : null}
      </div>
    </div>)
  }
}
