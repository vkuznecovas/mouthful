import { h, Component } from 'preact';
import style from './style';
import Thread from './thread';
import Login from './login';

const handleStateChange = (http, context, key) => {
		if(http.readyState == 4 && http.status == 200) {
			var stateChange = {loaded: true, authorized: true}
			stateChange[key] = JSON.parse(http.responseText)
			context.setState(stateChange)
		} else if (http.readyState == 4 && http.status == 401) {
			context.setState({authorized: false, loaded: false})
		} else {
			// TODO
		}
}

export default class Profile extends Component {
	constructor() {
		super();
		this.state = { threads: [], comments: [], authorized: false, loaded: false };
		this.loadThreads = this.loadThreads.bind(this);
		this.loadComments = this.loadComments.bind(this);
		this.loggedIn = this.loggedIn.bind(this);
	}

	loadThreads(context) {
		if (typeof window == "undefined") { return }
		
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/threads";
		http.open("GET", url, true);
		
		http.onreadystatechange = function() {
			handleStateChange(http, context, "threads")
		}
		http.send()
	}

	loadComments(context) {
		if (typeof window == "undefined") { return }
		
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/comments/all";
		http.open("GET", url, true);
		
		http.onreadystatechange = function() {
			handleStateChange(http, context, "comments")
		}
		
		http.send()
	}

	loggedIn() {
		this.setState({authorized: true})
	}

	// Note: `user` comes from the URL, courtesy of our router
	render({ user }) {
		if (!this.state.loaded) {
			this.loadThreads(this)
			this.loadComments(this)
		}
		if (!this.state.authorized) {
			return (<Login onLogin={this.loggedIn} />)
		}
		const threads = []
		for (var i =0; i< this.state.threads.length; i++){
			var comments = this.state.comments.filter(x => {
				console.log(x.ThreadId == this.state.threads[i].Id);
				return x.ThreadId == this.state.threads[i].Id;
			})
			if (comments.length == 0) {
				continue;
			}
			threads.push(<Thread thread={this.state.threads[i]} comments={comments}/>)
		}
		return (
			<div class={style.profile}>
				<div>
					{threads}
				</div>
			</div>	
		);
	}
}
