import { h, Component } from 'preact';
import style from './style';
import Thread from './thread';
import Login from './login';

const handleStateChange = (http, context, key) => {
	if (http.readyState == 4 && http.status == 200) {
		var stateChange = { loaded: true, authorized: true }
		stateChange[key] = JSON.parse(http.responseText)
		console.log('key', key);
		console.log(JSON.stringify(JSON.parse(http.responseText)));
		context.setState(stateChange)
	} else if (http.readyState == 4 && http.status == 401) {
		context.setState({ authorized: false, loaded: false })
	} else {
		// TODO
	}
}

export default class Panel extends Component {
	constructor() {
		super();
		this.state = { threads: [], comments: [], authorized: false, loaded: false, showPending: false };
		this.loadThreads = this.loadThreads.bind(this);
		this.loadComments = this.loadComments.bind(this);
		this.loggedIn = this.loggedIn.bind(this);
		this.showPending = this.showPending.bind(this);
		this.hidePending = this.hidePending.bind(this);
		this.reload = this.reload.bind(this);
	}

	showPending() {
		this.setState({ showPending: true })
	}

	hidePending() {
		this.setState({ showPending: false })
	}

	loadThreads(context) {
		if (typeof window == "undefined") { return }

		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/threads";
		http.open("GET", url, true);

		http.onreadystatechange = function () {
			handleStateChange(http, context, "threads")
		}
		http.send()
	}

	loadComments(context) {
		if (typeof window == "undefined") { return }

		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/comments/all";
		http.open("GET", url, true);

		http.onreadystatechange = function () {
			handleStateChange(http, context, "comments")
		}

		http.send()
	}

	loggedIn() {
		this.setState({ authorized: true })
	}

	reload() {
		this.setState({ loaded: false, threads: [], comments: [] })
	}

	render() {
		if (!this.state.loaded) {
			this.loadThreads(this)
			this.loadComments(this)
		}
		if (!this.state.authorized) {
			return (<Login onLogin={this.loggedIn} />)
		}
		// if we have no threads or comments
		if (!(this.state.threads && this.state.comments && this.state.threads.length && this.state.comments.length)) {
			return <div class={style.profile}>No comments yet!</div>
		}

		var threads = this.state.threads.map(t => {
			var comments = this.state.comments
			var c = comments.filter(x => {
				if (x.ThreadId == t.Id) {
					if (this.state.showPending) {
						return !x.Confirmed
					}
					return x.Confirmed;
				}
				return false;
			})
			if (c.length != 0) {
				return <Thread key={"___thread" + t.Id} thread={t} comments={c} reload={this.reload} />
			}
			return null;
		})
		console.log("threads", threads);

		return (
			<div class={style.profile}>
				<div class="buttons">
					<div onClick={this.showPending}>Pending</div>
					<div onClick={this.hidePending}>Verified</div>
				</div>
				<div>
					{threads.filter(x => x != null)}
				</div>
			</div>
		);
	}
}
