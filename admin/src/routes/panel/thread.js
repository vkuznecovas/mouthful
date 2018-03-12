import { h, Component } from 'preact';
import style from './style';


function formatDate(d) {
	var dd = new Date(d)
	return dd.toISOString().slice(0,19).replace("T", " ")
}

export default class Thread extends Component {
	constructor(props) {
		super(props);
		this.reload = this.reload.bind(this)
		this.deleteComment = this.deleteComment.bind(this)
		this.updateComment = this.updateComment.bind(this)
		this.undoDelete = this.undoDelete.bind(this)
	}
	reload() {
		this.props.reload()
	}
	deleteComment(commentId) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/v1/admin/comments";
		http.open("DELETE", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				var comments = context.props.comments.filter(x => x.Id != commentId)
				context.setState({comments})
				context.reload()
			} else if (http.readyState == 4 && http.status == 401) {
				context.reload()
			} else {
				context.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId }))
	}
	undoDelete(commentId) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/v1/admin/comments/restore";
		http.open("POST", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				var comments = context.props.comments.filter(x => x.Id != commentId)
				context.setState({comments})
				context.reload()
			} else if (http.readyState == 4 && http.status == 401) {
				context.reload()
			} else {
				context.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId }))
	}
	updateComment(commentId, body, author, confirmed) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/v1/admin/comments";
		http.open("PATCH", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				var comments = context.props.comments
				var comment = comments.filter(x => x.Id == commentId)
				if (comment) {
					comment[0].Body = body
					comment[0].Author = author
					comment[0].Confirmed = confirmed
				}
			} else if (http.readyState == 4 && http.status == 401) {
				context.reload()
			} else {
				context.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId, Body: body, Author: author, Confirmed: confirmed }))
	}
	

	render() {
		let comments = ""
		// this is a terrible hack for showing unconfirmeds
		if (this.props.comments.filter(comment => comment.ReplyTo == null && comment.DeletedAt == null).length == 0) {
			comments = this.props.comments.map(comment => {
				return <div class={style.comment}>
					<div  class={style.author}>By: {comment.Author} </div>
					<div class={style.date}>{formatDate(comment.CreatedAt)}</div>
					<div><textarea value={comment.Body}></textarea></div>
					<div class={style.buttons}>
						<div class={style.smallButton}  onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, comment.Confirmed)}>Update</div>
						{comment.DeletedAt == null ? <div class={style.smallButton}  onClick={() => this.deleteComment(comment.Id)}>Delete</div> : <div class={style.smallButton} onClick={() => this.undoDelete(comment.Id)}>Undo delete</div>}
						{comment.Confirmed ? "" : <div class={style.smallButton} onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, true)}>Confirm</div>}
					</div>
				</div>;
			})
		} else {
			comments = this.props.comments.filter(comment => comment.ReplyTo == null || comment.DeletedAt != null).map(comment => {
				var replies = this.props.comments.filter(x => x.ReplyTo === comment.Id).map(x => {
					return <div class={style.commentReply}>
						<div  class={style.author}>By: {x.Author}</div>
						<div class={style.date}>{formatDate(x.CreatedAt)}</div>
						<div><textarea value={x.Body}></textarea></div>
						<div class={style.buttons}>
							<div class={style.smallButton} onClick={() => this.updateComment(x.Id, x.Body, x.Author, x.Confirmed)}>Update</div>
							{x.DeletedAt == null ? <div class={style.smallButton} onClick={() => this.deleteComment(x.Id)}>Delete</div> : <div class={style.smallButton} onClick={() => this.undoDelete(x.Id)}>Undo delete</div>}
							{x.Confirmed ? "" : <div class={style.smallButton} onClick={() => this.updateComment(x.Id, x.Body, x.Author, true)}>Confirm</div> }
						</div>
					</div>
				});
				return <div class={style.comment}>
					<div class={style.author}>By: {comment.Author}</div>
					<div class={style.date}>{formatDate(comment.CreatedAt)}</div>
					<div><textarea value={comment.Body}></textarea></div>
					<div class={style.buttons}>
						<div class={style.smallButton} onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, comment.Confirmed)}>Update</div>
						{comment.DeletedAt == null ? <div class={style.smallButton} onClick={() => this.deleteComment(comment.Id)}>Delete</div> : <div class={style.smallButton} onClick={() => this.undoDelete(comment.Id)}>Undo delete</div>}
						{comment.Confirmed ? "" : <div class={style.smallButton} onClick={() => this.updateComment(comment.Id, comment.Body, comment.Author, true)}>Confirm</div> }
					</div>
					<div style="margin-left:30px">
						{replies}
					</div>
				</div>;
			})
		}
		const fullThread = comments.length > 0 ? (<div class="thread">
		<h2>{this.props.thread.Path}</h2>
		<div class="comments">
			{comments}
		</div>
	</div>) : null
		return fullThread
		
	}
}
