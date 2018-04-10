import { h, Component } from 'preact';
import style from './style';
import Thread from './thread';
import Login from './login';

const handleStateChange = (http, context, key) => {
	if (http.readyState == 4 && http.status == 200) {
		var stateChange = { loaded: true, authorized: true }
		stateChange[key] = JSON.parse(http.responseText)
		context.setState(stateChange)
	} else if (http.readyState == 4 && http.status == 401) {
		context.setState({ authorized: false, loaded: true })
	}
}

export default class Panel extends Component {
	constructor() {
		super();
		this.state = { threads: [], comments: [], authorized: false, loaded: false, showPending: true, showDeleted: false };
		this.loadThreads = this.loadThreads.bind(this);
		this.loadComments = this.loadComments.bind(this);
		this.loggedIn = this.loggedIn.bind(this);
		this.showPending = this.showPending.bind(this);
		this.hidePending = this.hidePending.bind(this);
		this.reload = this.reload.bind(this);
		this.showDeleted = this.showDeleted.bind(this);
		this.updateComment = this.updateComment.bind(this);
	}

	showPending() {
		this.setState({ showPending: true })
		this.setState({ showDeleted: false})
	}

	hidePending() {
		this.setState({ showPending: false })
		this.setState({ showDeleted: false })
	}

	showDeleted() {
		this.setState({ showPending: false })
		this.setState({ showDeleted: true })
	}

	loadThreads(context) {
		if (typeof window == "undefined") { return }

		var http = new XMLHttpRequest();
		var url = window.location.origin + "/v1/admin/threads";
		http.open("GET", url, true);

		http.onreadystatechange = function () {
			handleStateChange(http, context, "threads")
		}
		http.send()
	}

	loadComments(context) {
		if (typeof window == "undefined") { return }

		var http = new XMLHttpRequest();
		var url = window.location.origin + "/v1/admin/comments/all";
		http.open("GET", url, true);

		http.onreadystatechange = function () {
			handleStateChange(http, context, "comments")
		}

		http.send()
	}

	updateComment(commentId, body, author, confirmed) {
		if (typeof window == "undefined") { return }
		var http = new XMLHttpRequest();
		var url = window.location.origin + "/v1/admin/comments";
		http.open("PATCH", url, true);
		var context = this;
		http.onreadystatechange = function () {
			if (http.readyState == 4 && http.status == 204) {
				context.reload()
			} else if (http.readyState == 4 && http.status == 401) {
				context.reload()
			} else {
				context.reload()
			}
		}
		http.send(JSON.stringify({ CommentId: commentId, Body: body, Author: author, Confirmed: confirmed }))
	}
	

	loggedIn() {
		this.setState({ authorized: true })
		this.setState({ loaded: false })
	}

	reload() {
		this.setState({ loaded: false, threads: [], comments: [] })
	}
	componentWillMount() {
		this.loggedIn();
	}
	render() {
		if (!this.state.authorized) {
			return (<div class={style.mouthful_container}>
			<Login onLogin={this.loggedIn} />
			</div>)
		}
		if (!this.state.loaded) {
			this.loadThreads(this)
			this.loadComments(this)
		}
		
		// if we have no threads or comments
		if (!(this.state.threads && this.state.comments && this.state.threads.length && this.state.comments.length)) {
			return <div class={style.mouthful_container}><div class={style.mouthful_login}>No comments yet!</div></div>
		}
		var threads = this.state.threads.map(t => {
			var comments = this.state.comments
			const pendingFilter = x => {
				return !x.Confirmed && x.DeletedAt == null
			}
			const showAll = x => {
				return x.DeletedAt == null
			}
			const deletedFilter = x => {
				return x.DeletedAt != null
			}
			
			var c = comments.filter(x => {
				return x.ThreadId == t.Id
			})
			let filter = this.state.showPending ? pendingFilter : showAll
			filter = this.state.showDeleted ? deletedFilter : filter

			c = c.filter(filter)
			if (c.length != 0) {
				return <Thread key={"___thread" + t.Id} thread={t} comments={c} reload={this.reload} updateComment={this.updateComment} />
			}
			return null;
		})
		var resultDiv = threads.filter(x => x != null).length > 0 ? threads : <div class={style.nothing}>Nothing to display</div>
		return (
			<div class={style.mouthful_container}>
			<div class={style.mouthful_wrapper}>
				<div class={style.mouthful_buttons}>
					<div class={this.state.showPending ? style.mouthful_buttonActive : style.mouthful_button} onClick={this.showPending}>Show unconfirmed</div>
					<div class={this.state.showPending == false && this.state.showDeleted == false ? style.mouthful_buttonActive : style.mouthful_button} onClick={this.hidePending}>Show all</div>
					<div class={this.state.showDeleted ? style.mouthful_buttonActive : style.mouthful_button  } onClick={this.showDeleted}>Show deleted</div>
				</div>
				<div>
					{resultDiv}
				</div>
			</div>
			</div>
		);
	}
}
