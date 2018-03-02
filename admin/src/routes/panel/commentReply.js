import { h, Component } from 'preact';
import style from './style';


export default class CommentReply extends Component {
	constructor(props) {
		super(props);
		this.state = { comment: comment, replies: props.replies };
		this.reload = props.reload.bind(this);
		this.deleteComment = this.deleteComment.bind(this)
		this.updateComment = this.updateComment.bind(this)
	}
	deleteComment(commentId) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/comments";
		http.open("DELETE", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				var comments = context.state.comments.filter(x => x.Id != commentId)
				context.setState({comments})
				this.reload()
			} else if (http.readyState == 4 && http.status == 401) {
				console.log("TODO");
				this.reload()
			} else {
				console.log("TODO");
				this.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId }))
	}
	updateComment(commentId, body, author, confirmed) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/comments";
		http.open("PATCH", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				var comments = context.state.comments
				var comment = comments.filter(x => x.Id == commentId)
				if (comment) {
					comment[0].Body = body
					comment[0].Author = author
					comment[0].Confirmed = confirmed
				}
			} else if (http.readyState == 4 && http.status == 401) {
				console.log("TODO");
				this.reload()
			} else {
				console.log("TODO");
				this.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId, Body: body, Author: author, Confirmed: confirmed }))
	}
	

	render() {
		const comments = this.state.comments.filter(comment => !comment.ReplyTo).map(comment => {
			var replies = this.state.comments.filter(x => x.ReplyTo === comment.Id).map(x => {
				return <div class="comment-reply">
					<h3>By: {x.Author} <span>at {x.CreatedAt}</span></h3>
					<textarea value={x.Body}></textarea>
					<div>
						<div onClick={() => this.updateComment(x.Id, x.Body, x.Author, x.Confirmed)}>Update</div>
						<div onClick={() => this.deleteComment(x.Id)}>Delete</div>
						{x.Confirmed ? "" : <div onClick={() => this.updateComment(x.Id, x.Body, x.Author, true)}>Confirm</div>}
					</div>
				</div>
			});
			return <div class="comment">
				<h3>By: {comment.Author} <span>at {comment.CreatedAt}</span></h3>
				<textarea value={comment.Body}></textarea>
				<div>
				<div onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, comment.Confirmed)}>Update</div>
				<div onClick={() => this.deleteComment(comment.Id)}>Delete</div>
				{comment.Confirmed ? "" : <div onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, true)}>Confirm</div>}
			</div>
				<div style="margin-left:30px">
					{replies}
				</div>
			</div>;
		})
		return (
			<div class="thread">
				<h2>{this.state.thread.Path}</h2>
				<div class="comments">
					{comments}
				</div>
			</div>
		);
	}
}
