import { h, Component } from 'preact';
import style from './style';
import { route } from 'preact-router';


export default class Login extends Component {
	constructor(props) {
        super(props);
		this.state = { value: '' };
        this.onLogin = props.onLogin;
        
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}

	handleSubmit(event) {
		var context = this;
        
		event.preventDefault();		
		var http = new XMLHttpRequest();
		var url = "http://localhost:7777/admin/login";
		http.open("POST", url, true);
		
		//Send the proper header information along with the request
		http.setRequestHeader("Content-type", "application/json");
		http.onreadystatechange = function() {//Call a function when the state changes.
			if(http.readyState == 4 && http.status == 204) {
                context.onLogin();
			} else {
				// TODO
			}
		}
		http.send(JSON.stringify({password: context.state.value}));
	}
	render() {
		return (
			<div class={style.login}>
				<form onSubmit={this.handleSubmit}>
					<label>
						Password:
					<input type="password" value={this.state.value} onChange={this.handleChange} />
					</label>
					<input type="submit" value="Submit" />
				</form>
			</div>
		);
	}
}
